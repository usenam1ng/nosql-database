package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// Функция для генерации случайной строки
func generateRandomString() string {
	length := 7
	rand.Seed(time.Now().UnixNano())

	// Задаем символы, из которых будет формироваться строка
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdrfghijklmnjpqrstuvwxyz0123456789"

	str := make([]byte, length)

	// Генерируем случайную строку
	for i := 0; i < length; i++ {
		str[i] = charset[rand.Intn(len(charset))]
	}

	return string(str)
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method == http.MethodPost {
		baseurl := r.FormValue("url")
		newurl := generateRandomString()

		recuesttodb := "--file filellt.data --query 'HSET attttt " + newurl + " " + baseurl + "'"

		conn, err := net.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		if _, err = conn.Write([]byte(recuesttodb)); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(w, "Your short URL: 127.0.0.1:8080/"+newurl)

	} else if r.Method == http.MethodGet {
		newurl := r.URL.Path[1:]

		recuesttodb := "--file filellt.data --query 'HGET attttt " + newurl + " " + "tt" + "'"

		conn, err := net.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		if _, err = conn.Write([]byte(recuesttodb)); err != nil {
			fmt.Println(err)
			return
		}

		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}

		longurl := string(buffer[:n])

		http.Redirect(w, r, longurl, http.StatusSeeOther)

	}
}

func main() {
	http.HandleFunc("/", handle)             // Устанавливаем роутер
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
