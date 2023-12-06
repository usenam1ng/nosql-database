package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

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
	var bs []JSONE
	ans := make(map[string]interface{})

	file, err := os.ReadFile("stat.json")
	if err != nil {
		fmt.Println("нету файла")
		return
	}

	if len(file) == 0 {
		fmt.Println("пустой файл, соберите статистику")
		return
	}

	err = json.Unmarshal(file, &bs)
	if err != nil {
		fmt.Println("не анмаршалит")
		return
	}

	for _, entry := range bs {
		if entry.SourceIP != "" && entry.Time != "" {
			ans[entry.SourceIP] = make(map[string]interface{})
			var minTime, maxTime time.Time
			minTime, err := time.Parse("2006-01-02 15:04:05", entry.Time)
			if err != nil {
				fmt.Println("Ошибка парсинга даты 1:", err)
				return
			}
			for _, i := range bs {
				if i.Time != "" {
					data, err := time.Parse("2006-01-02 15:04:05", i.Time)
					if err != nil {
						fmt.Println("Ошибка парсинга даты 2:", err)
						return
					}
					if data.Before(minTime) {
						minTime = data
					}
					if data.After(maxTime) {
						maxTime = data
					}
				}
			}
			mTime := minTime.Format("2006-01-02 15:04:05")
			mxTime := maxTime.Format("2006-01-02 15:04:05")
			ax := mTime + " - " + mxTime

			ans[ax] = make(map[string]interface{})

			for _, i := range bs {
				if i.ShortURL != " " && i.ShortURL != "" && i.ShortURL != "\n" {
					ans[ax].(map[string]interface{})[i.ShortURL] = i.Count
				}
			}
		}
	}

	jsonData, err := json.MarshalIndent(ans, "", "  ")
	if err != nil {
		fmt.Println("не маршалит")
		return
	}

	err = os.WriteFile("report.json", jsonData, 0644)
	if err != nil {
		fmt.Println("не пишет в файл")
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
		fmt.Println("нету файла")
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

	report()
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
