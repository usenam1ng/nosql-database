package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type connectionReport struct {
	ShortUrl string `json:"shortURL"`
	OutLink  string `json:"outLink"`
	Host     string `json:"originHost"`
}

type JSONE struct {
	ID       int    `json:"id"`
	PID      int    `json:"pid"`
	URL      string `json:"url"`
	ShortURL string `json:"shortURL"`
	SourceIP string `json:"sourceIP"`
	Time     string `json:"time"`
	Count    int    `json:"count"`
}

type Payload struct {
	Fullstat []string `json:"Fullstat"`
}

func pid(conn []JSONE, url string) int {
	PID := 0
	for _, connect := range conn {
		if connect.URL == url {
			PID = connect.ID
		}
	}
	return PID
}
func id(conn []JSONE) int {
	maxID := 0
	for _, connect := range conn {
		if connect.ID > maxID {
			maxID = connect.ID
		}
	}
	return maxID + 1
}

func par(conn []JSONE, url string) bool {
	for _, connect := range conn {
		if connect.URL == url {
			return false
		}
	}
	return true
}

func count(conn []JSONE, url string) {
	for index := range conn {
		if conn[index].URL == url {
			conn[index].Count++
			return
		}
	}
}

func statConnections(url, shortURL, ip string) {
	var connect []JSONE

	parent_conn := JSONE{
		URL:      url,
		ShortURL: shortURL,
		Count:    1,
	}

	new_conn := JSONE{
		SourceIP: ip,
		Time:     time.Now().Format(""),
		Count:    1,
	}

	file, err := os.ReadFile("stat.json")
	if err != nil {
		fmt.Println("no file")
		return
	}

	if len(file) == 0 {
		connect = []JSONE{}
	}

	err = json.Unmarshal(file, &connect)
	if err != nil {
		fmt.Println("не анмаршалит")
		return
	}

	if connect == nil {
		connect = []JSONE{}
	}

	parent_conn.ID = id(connect)

	if par(connect, parent_conn.URL) == true {
		connect = append(connect, parent_conn)
	} else {
		count(connect, parent_conn.URL)
	}

	new_conn.ID = id(connect)
	new_conn.PID = pid(connect, url)
	connect = append(connect, new_conn)

	jsonData, err := json.MarshalIndent(connect, "", "  ")
	if err != nil {
		fmt.Println("не маршалит")
		return
	}

	err = os.WriteFile("stat.json", jsonData, 0644)
	if err != nil {
		fmt.Println("не пишет в файл")
		return
	}
}

func handle(conn net.Conn) {
	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	line := string(buffer[:n])
	elements := strings.Split(line, " ")

	statConnections(elements[0], elements[1], elements[2])

	return

}

func main() {
	listener, err := net.Listen("tcp", ":6575")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err.Error())
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err.Error())
			continue
		}
		go handle(conn)
	}
}
