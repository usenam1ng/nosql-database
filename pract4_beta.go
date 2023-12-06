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
	ShortURL string `json:"ShortUrl"`
	URL      string `json:"URL"`
	SourceIP string `json:"sourceIP"`
	Time     string `json:"time"`
	Count    int    `json:"count"`
}

type Payload struct {
	Fullstat []string `json:"Fullstat"`
}

func manipulateJSONE(conn []JSONE, url string) (int, int, bool) {
	PID := 0
	maxID := 0
	indexURL := -1

	for index, connect := range conn {
		if connect.URL == url {
			PID = connect.ID
			indexURL = index
			conn[index].Count++
		}
		if connect.ID > maxID {
			maxID = connect.ID
		}
	}

	if indexURL == -1 {
		return maxID + 1, PID, true
	}

	return maxID + 1, PID, false
}

func report() {
	var connect []JSONE
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
		Time:     time.Now().Format("2006-01-02 15:04:05"),
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

	newPID, existingPID, isNew := manipulateJSONE(connect, parent_conn.URL)

	parent_conn.ID = newPID

	if isNew == true {
		connect = append(connect, parent_conn)
	} else {
		for index := range connect {
			if connect[index].URL == url {
				connect[index].Count++
				return
			}
		}
	}

	newPID, existingPID, isNew = manipulateJSONE(connect, parent_conn.URL)

	new_conn.ID = newPID
	new_conn.PID = existingPID
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

	fmt.Println("Если хотите получить отчет введите - 1, если нет - 0")
	var a int
	fmt.Scan(&a)
	if a == 1 {
		report()
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
