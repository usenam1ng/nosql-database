package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
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

func gen_pid(conn []JSONE, url string) int {
	PID := 0
	for _, connect := range conn {
		if connect.URL == url {
			PID = connect.ID
		}
	}
	return PID
}
func gen_unpid(conn []JSONE) int {
	maxID := 0
	for _, connect := range conn {
		if connect.ID > maxID {
			maxID = connect.ID
		}
	}
	return maxID + 1
}

func unique_par(conn []JSONE, url string) bool {
	for _, connect := range conn {
		if connect.URL == url {
			return false
		}
	}
	return true
}

func par_count(conn []JSONE, url string) {
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
		Time:     time.Now().Format("2999-01-02 00:00"),
		Count:    1,
	}

	file, err := os.ReadFile("statА у а у.json")
	if err != nil {
		fmt.Println("no file")
		return
	}

	if len(file) == 0 {
		fmt.Println("len file 0")
		return
	}

	err = json.Unmarshal(file, &connect)
	if err != nil {
		fmt.Println("не анмаршалит")
		return
	}

	if connect == nil {
		connect = []JSONE{}
	}

	parent_conn.ID = gen_unpid(connect)

	if unique_par(connect, parent_conn.URL) == true {
		connect = append(connect, parent_conn)
	} else {
		par_count(connect, parent_conn.URL)
	}

	new_conn.ID = gen_unpid(connect)
	new_conn.PID = gen_pid(connect, url)
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

func handle(w http.ResponseWriter, r *http.Request) {
	var report connectionReport

	err := json.NewDecoder(r.Body).Decode(&report)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statConnections(report.OutLink, report.ShortUrl, report.Host)

	return

}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handle)

	log.Fatal(http.ListenAndServe(":6565", nil))
}
