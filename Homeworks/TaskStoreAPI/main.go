/*
1. Реализовать TaskStoreAPI с помощью GorillaMux
2. Реализовать в качестве хранилища SQLite
3. Добавить реализацию endpoints для '/tags/' и '/due/'
4. Важное замечание про тип time.Time и SQLite.
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"./taskstore/"

	"github.com/gorilla/mux"
)

func main() {
	zLog("help run:   SERVERPORT=8888 go run . \n")
	zPort := os.Getenv("SERVERPORT")
	if zPort == "" {
		zPort = "1234"
	}
	zHostPort := "0.0.0.0:" + zPort
	zRouter := mux.NewRouter()
	zRouter.StrictSlash(true)
	zRouter.HandleFunc("/info", zInfo).Methods("GET")
	zRouter.HandleFunc("/task/", zGetTasksAll).Methods("GET")

	zLog("Started at " + zHostPort + "\n")
	log.Fatal(http.ListenAndServe(zHostPort, nil)) // Запускаем web-сервер в режиме "слушания"
}

type ztDB struct {
	zStore *taskStore
}

func zDBinit() {
	zDB := taskstore.New()
	return &zDB{zStore: zStore}
}

func zGetTasksAll(w http.ResponseWriter, r *http.Request) {

}

func zInfo(w http.ResponseWriter, r *http.Request) {
	zLog("zInfo")
	fmt.Fprintf(w, "\n<h1>TaskStoreAPI</h1>")
	fmt.Fprintf(w, "\n<p><a href='/info'   target='_blank'>/info</a>       // Информация об API</p>")
	fmt.Fprintf(w, "\n<p><a href='/task'   target='_blank'>/tasks_all</a>  // возвращает все задачи</p>")
	zLog(".\n")
}

func zLog(txt string) {
	vCurrTime := time.Now()
	if txt != ".\n" {
		fmt.Print(" ; ", vCurrTime.Format("15:04:05.000"), " ")
	}
	fmt.Print(txt)
}
