package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// w - responseWriter (куда писать ответ)
// r - request (откуда брать запрос)
// Функция-обработчик(Handler)
func GetGreet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hi! I'm new web-server</h1>")
}

func zPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p>ping back</p>")
}

func main() {
	http.HandleFunc("/", GetGreet) // Если придет запрос на адрес "/", то вызывай GetGreet
	http.HandleFunc("/ping", zPing)
	// Запуск сервера в консоли: SERVERPORT=5000 go run .
	fmt.Println("help run:  SERVERPORT=5000 go run . ")
	fmt.Println("0.0.0.0:" + os.Getenv("SERVERPORT"))
	log.Fatal(http.ListenAndServe("0.0.0.0:"+os.Getenv("SERVERPORT"), nil)) // Запускаем web-сервер в режиме "слушания"
}
