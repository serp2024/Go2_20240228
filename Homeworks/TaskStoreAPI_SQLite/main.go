/*
1. Реализовать TaskStoreAPI с помощью GorillaMux
2. Реализовать в качестве хранилища SQLite
3. Добавить реализацию endpoints для '/tags/' и '/due/'
4. Важное замечание про тип time.Time и SQLite.
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gorilla/mux"
)

type Task struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
	Tags string `json:"tags"`
	Due  string `json:"due"`
	//DueN string `json:"dueN"`
}

type ResponseId struct {
	Id int `json:"id"`
}

var Zdb *sql.DB
var Err error

func zDBinit() {

	os.Remove("./todo.db")
	Zdb, Err = sql.Open("sqlite", "./todo.db")
	if Err != nil {
		log.Fatal(Err)
	}
	//defer Zdb.Close()

	zSQL := `create table if not exists task(id INTEGER not null primary key AUTOINCREMENT, txt text, tags text, due_txt text, due_int integer);
	         delete from task;`
	_, err := Zdb.Exec(zSQL)
	if err != nil {
		zLog("Error: " + err.Error())
		panic(err)
	}
}

func zTaskGetAll(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetAll ")
	zLog("p0 ")

	rows, err := Zdb.Query("select id,txt,tags,due_txt from task")
	if err != nil {
		log.Fatal(err)
	}

	zLog("p1 ")

	var zResArr []Task
	for rows.Next() {
		var t1 Task
		rows.Scan(&t1.Id, &t1.Text, &t1.Tags, &t1.Due)
		zResArr = append(zResArr, t1)
	}
	zRes, _ := json.Marshal(zResArr)
	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}

func zTaskCreate(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskCreate ")

	var t1 Task

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&t1)

	zTime, _ := time.Parse(time.RFC3339, t1.Due)
	zDueN := zTime.Unix()

	result, err := Zdb.Exec("INSERT INTO task(txt, tags, due_txt, due_int) VALUES(?, ?, ?, ?) ", t1.Text, t1.Tags, t1.Due, zDueN)
	if err != nil {
		zErrorResponse(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
	id, err := result.LastInsertId()
	if err != nil {
		zErrorResponse(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
	zLog("LastInsertId = " + strconv.Itoa(int(id)))
	zRes, err := json.Marshal(ResponseId{Id: int(id)})
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}

func zTaskDelete(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskDelete ")

	zVars := mux.Vars(r)
	zId, err := strconv.Atoi(zVars["id"])
	if err != nil {
		zErrorResponse(w, "Invalid ID ("+zVars["id"]+" must be int)", http.StatusBadRequest)
	}
	Zdb.Query("delete from task where id=?", zId)
	zLog(".\n")
}

func zTaskDeleteAll(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskDeleteAll ")
	Zdb.Query("delete from task")
	zLog(".\n")
}

func zTaskGetByTag(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetByTag ")
	zVars := mux.Vars(r)

	// про безопасноть от SQL инекций в учебном ТЗ не было сказано ;)
	// для поиска по полному тэгу (при условии что тэги разделены пробелами)
	//zSQL := "select id,txt,tags,due_txt from task where ' '||tags||' ' like '% " + zVars["name"] + " %'"
	// или поиск попроще
	zSQL := "select id,txt,tags,due_txt from task where tags like '%" + zVars["name"] + "%'"

	zLog("\n" + zSQL + "\n")
	rows, err := Zdb.Query(zSQL)
	if err != nil {
		log.Fatal(err)
	}

	var zResArr []Task
	for rows.Next() {
		var t1 Task
		rows.Scan(&t1.Id, &t1.Text, &t1.Tags, &t1.Due)
		zResArr = append(zResArr, t1)
	}
	zRes, _ := json.Marshal(zResArr)
	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}

func zTaskGetByYMD_t1t2(w http.ResponseWriter, t1 int64, t2 int64) {
	zLog("zTaskGetByYMD_t1t2 ")

	zSQL := "select id,txt,tags,due_txt from task where due_int between ? and ?"

	rows, err := Zdb.Query(zSQL, t1, t2)
	if err != nil {
		log.Fatal(err)
	}

	var zResArr []Task
	for rows.Next() {
		var t1 Task
		rows.Scan(&t1.Id, &t1.Text, &t1.Tags, &t1.Due)
		zResArr = append(zResArr, t1)
	}
	zRes, _ := json.Marshal(zResArr)
	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}
func zTaskGetByYMD(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetByYMD ")
	zVars := mux.Vars(r)

	//проверка на корректность даты
	zYMD := zVars["yyyy"] + "-" + zVars["mm"] + "-" + zVars["dd"] + "T00:00:00.000Z" // 2024-03-04T20:21:22.123Z
	t, err := time.Parse(time.RFC3339, zYMD)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	zLog("check date: " + t.String())

	t, _ = time.Parse(time.RFC3339, zVars["yyyy"]+"-"+zVars["mm"]+"-"+zVars["dd"]+"T00:00:00.000Z")
	t1 := t.Unix()
	t, _ = time.Parse(time.RFC3339, zVars["yyyy"]+"-"+zVars["mm"]+"-"+zVars["dd"]+"T23:59:59.000Z")
	t2 := t.Unix()

	zTaskGetByYMD_t1t2(w, t1, t2)
	zLog(".\n")
}

func zTaskGetByYM(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetByYMD ")
	zVars := mux.Vars(r)

	//проверка на корректность даты
	zYMD := zVars["yyyy"] + "-" + zVars["mm"] + "-01" + "T00:00:00.000Z" // 2024-03-04T20:21:22.123Z
	t, err := time.Parse(time.RFC3339, zYMD)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	zLog("check date: " + t.String())

	t, _ = time.Parse(time.RFC3339, zVars["yyyy"]+"-"+zVars["mm"]+"-01"+"T00:00:00.000Z")
	t1 := t.Unix()
	t2 := t.AddDate(0, 1, 0).Unix()

	zTaskGetByYMD_t1t2(w, t1, t2)
	zLog(".\n")
}

func zTaskGetByY(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetByYMD ")
	zVars := mux.Vars(r)

	//проверка на корректность даты
	zYMD := zVars["yyyy"] + "-01-01" + "T00:00:00.000Z" // 2024-03-04T20:21:22.123Z
	t, err := time.Parse(time.RFC3339, zYMD)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	zLog("check date: " + t.String())

	t, _ = time.Parse(time.RFC3339, zVars["yyyy"]+"-01-01"+"T00:00:00.000Z")
	t1 := t.Unix()
	t2 := t.AddDate(1, 0, 0).Unix()

	zTaskGetByYMD_t1t2(w, t1, t2)
	zLog(".\n")
}

func main() {
	zLog("help run:   SERVERPORT=5555 go run . \n")
	zPort := os.Getenv("SERVERPORT")
	if zPort == "" {
		zPort = "5555"
	}
	zHostPort := "0.0.0.0:" + zPort

	zDBinit()

	zRouter := mux.NewRouter()
	// zRouter.StrictSlash(true)
	zRouter.HandleFunc("/", zInfo).Methods("GET")
	zRouter.HandleFunc("/info", zInfo).Methods("GET")                                                    //1
	zRouter.HandleFunc("/task/", zTaskCreate).Methods("POST")                                            //2
	zRouter.HandleFunc("/task/", zTaskGetAll).Methods("GET")                                             //3
	zRouter.HandleFunc("/task/{id:[0-9]+}", zTaskDelete).Methods("DELETE")                               //4
	zRouter.HandleFunc("/task", zTaskDeleteAll).Methods("DELETE")                                        //5
	zRouter.HandleFunc("/tag/{name}", zTaskGetByTag).Methods("GET")                                      //6
	zRouter.HandleFunc("/due/{yyyy:[0-9]{4}}/{mm:[0-9]{2}}/{dd:[0-9]{2}}", zTaskGetByYMD).Methods("GET") //7
	zRouter.HandleFunc("/due/{yyyy:[0-9]{4}}/{mm:[0-9]{2}}", zTaskGetByYM).Methods("GET")                //8
	zRouter.HandleFunc("/due/{yyyy:[0-9]{4}}", zTaskGetByY).Methods("GET")                               //9

	zLog("Started at " + zHostPort + "\n")
	log.Fatal(http.ListenAndServe(zHostPort, zRouter)) // Запускаем web-сервер в режиме "слушания"

}

func zInfo(w http.ResponseWriter, r *http.Request) {
	zLog("zInfo")
	fmt.Fprintf(w, "\n<h1>TaskStoreAPI using SQLite</h1>")
	fmt.Fprintf(w, "\n<p>1) <a href='/info'   target='_blank'>/info</a>       // Информация об API</p>")
	fmt.Fprintf(w, `  <p>2) /task/ // Сохранить новую задачу (тэги разделять пробелом, запустить в консоли) // curl --header "Content-Type: application/json" --request POST --data '{"text":"test_text","tags":"tag1 tag2","due":"2024-03-04T20:21:22.123Z"}' http://127.0.0.1:5555/task/  </p>`)
	fmt.Fprintf(w, "\n<p>3) <a href='/task/'   target='_blank'>/tasks_all</a>  // возвращает все задачи // curl --request GET http://127.0.0.1:5555/task/  </p>")
	fmt.Fprintf(w, `  <p>4) /task/1 // Удалить задачу с номером 1 // curl --request DELETE   http://127.0.0.1:5555/task/1  </p>`)
	fmt.Fprintf(w, `  <p>5) /task // Удалить все задачи // curl --request DELETE   http://127.0.0.1:5555/task  </p>`)
	fmt.Fprintf(w, "\n<p>6) <a href='/tag/tag1'   target='_blank'>/tag/tag1</a>  // возвращает все задачи c тэгом tag1                       // curl --request GET http://127.0.0.1:5555/tag/tag1  </p>")
	fmt.Fprintf(w, "\n<p>7) <a href='/due/2024/03/04'   target='_blank'>/due/2024/03/04</a>    // возвращает все задачи с датой 2024-03-05   // curl --request GET http://127.0.0.1:5555/due/2024/03/04  </p>")
	fmt.Fprintf(w, "\n<p>7) <a href='/due/2024/03'      target='_blank'>/due/2024/03</a>       // возвращает все задачи с месяцем 2024-03    // curl --request GET http://127.0.0.1:5555/due/2024/03  </p>")
	fmt.Fprintf(w, "\n<p>7) <a href='/due/2024'         target='_blank'>/due/2024</a>          // возвращает все задачи с годом 2024         // curl --request GET http://127.0.0.1:5555/due/2024  </p>")
	zLog(".\n")
}

func zLog(txt string) {
	vCurrTime := time.Now()
	if txt != ".\n" {
		fmt.Print(" ; ", vCurrTime.Format("15:04:05.000"), " ")
	}
	fmt.Print(txt)
}

func zErrorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	zLog("zErrorResponse " + message + " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
