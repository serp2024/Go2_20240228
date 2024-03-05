/*
1. Реализовать TaskStoreAPI с помощью GorillaMux
2. Реализовать в качестве хранилища SQLite
3. Добавить реализацию endpoints для '/tags/' и '/due/'
4. Важное замечание про тип time.Time и SQLite.
*/

package main

import (
	"TaskStoreAPI/taskstore"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ztDB struct {
	zStore *taskstore.TaskStore
}

func zDBinit() *ztDB {
	zDB := taskstore.New()
	return &ztDB{zStore: zDB}
}

func (zDB ztDB) zTaskGetAll(w http.ResponseWriter, r *http.Request) {
	zLog("zGetTasksAll ")
 
	allTasks := zDB.zStore.GetAllTasks() 
	js, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	zLog(".\n")
}

func (zDB ztDB) zTaskCreate(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskCreate ")

	type RequestTask struct {
		Text string   `json:"text"`
		Tags []string `json:"tags"`
		Due  string   `json:"due"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	var zReqData RequestTask

	zContentType := r.Header.Get("Content-Type")
	zMediaType, _, err := mime.ParseMediaType(zContentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if zMediaType != "application/json" {
		http.Error(w, "POST data as  application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	zBody, _ := io.ReadAll(r.Body)
	//	zLog("zBody=" + string(zBody) + "; ")
	if err := json.Unmarshal([]byte(zBody), &zReqData); err != nil {
		zErrorResponse(w, err.Error(), http.StatusBadRequest)
		// fmt.Println("Error: ", err)
	}

	//	fmt.Print("zReqData=", zReqData, "; ")

	//zReqData.DisallowUnknownFields()
	// var zRT RequestTask
	// if err := zReqData.Decode(&zRT); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	zTime, err := time.Parse(time.RFC3339, zReqData.Due)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	//zLog("time1=" + zReqData.Due + "; ")
	//zLog("time2=" + zTime.String() + "; ")

	zTaskId := zDB.zStore.CreateTask(zReqData.Text, zReqData.Tags, zTime)
	zRes, err := json.Marshal(ResponseId{Id: zTaskId})
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusInternalServerError)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}

func (zDB ztDB) zTaskDelete(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskDelete ")

	zVars := mux.Vars(r)
	zId, err := strconv.Atoi(zVars["id"])
	if err != nil {
		zErrorResponse(w, "Invalid ID ("+zVars["id"]+" must be int)", http.StatusBadRequest)
	}
	err = zDB.zStore.DeleteTask(zId)
	if err != nil {
		zErrorResponse(w, "Error: id="+strconv.Itoa(zId)+" not found", 404)
	}

	zLog(".\n")
}

func (zDB ztDB) zTaskDeleteAll(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskDeleteAll ")

	err := zDB.zStore.DeleteAllTasks()
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	zLog(".\n")
}

func (zDB ztDB) zTaskGetByTag(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetByTag ")

	zVars := mux.Vars(r)
	zTagName := zVars["name"]

	zTaskArr := zDB.zStore.GetTaskByTag(zTagName)

	zRes, err := json.Marshal(zTaskArr)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}

func (zDB ztDB) zTaskGetByYMD(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetByYMD ")

	zVars := mux.Vars(r)
	zYYYY, err := strconv.Atoi(zVars["yyyy"])
	if err != nil {
		zErrorResponse(w, "Wrong year", http.StatusBadRequest)
		return
	}
	zMM, err := strconv.Atoi(zVars["mm"])
	if err != nil {
		zErrorResponse(w, "Wrong month", http.StatusBadRequest)
		return
	}
	zDD, err := strconv.Atoi(zVars["dd"])
	if err != nil {
		zErrorResponse(w, "Wrong day", http.StatusBadRequest)
		return
	}
	//проверка на корректность даты
	zYMD := zVars["yyyy"] + "-" + zVars["mm"] + "-" + zVars["dd"] + "T00:00:00.000Z" // 2024-03-04T20:21:22.123Z
	t, err := time.Parse(time.RFC3339, zYMD)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	zLog("check date: " + t.String())

	zTaskArr := zDB.zStore.GetTaskByYMD(zYYYY, zMM, zDD)

	zRes, err := json.Marshal(zTaskArr)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}

func (zDB ztDB) zTaskGetByYM(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetByYM ")

	zVars := mux.Vars(r)
	zYYYY, err := strconv.Atoi(zVars["yyyy"])
	if err != nil {
		zErrorResponse(w, "Wrong year", http.StatusBadRequest)
	}
	zMM, err := strconv.Atoi(zVars["mm"])
	if err != nil {
		zErrorResponse(w, "Wrong month", http.StatusBadRequest)
	}

	//проверка на корректность даты
	zYMD := zVars["yyyy"] + "-" + zVars["mm"] + "-01" + "T00:00:00.000Z" // 2024-03-04T20:21:22.123Z
	t, err := time.Parse(time.RFC3339, zYMD)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	zLog("check date: " + t.String())

	zTaskArr := zDB.zStore.GetTaskByYM(zYYYY, zMM)

	zRes, err := json.Marshal(zTaskArr)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}

func (zDB ztDB) zTaskGetByY(w http.ResponseWriter, r *http.Request) {
	zLog("zTaskGetByYM ")

	zVars := mux.Vars(r)
	zYYYY, err := strconv.Atoi(zVars["yyyy"])
	if err != nil {
		zErrorResponse(w, "Wrong year", http.StatusBadRequest)
	}

	zTaskArr := zDB.zStore.GetTaskByY(zYYYY)

	zRes, err := json.Marshal(zTaskArr)
	if err != nil {
		zErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(zRes)

	zLog(".\n")
}

func main() {
	zLog("help run:   SERVERPORT=5555 go run . \n")
	zPort := os.Getenv("SERVERPORT")
	if zPort == "" {
		zPort = "5555"
	}
	zHostPort := "0.0.0.0:" + zPort
	zDB := zDBinit()

	zDBinit()

	zRouter := mux.NewRouter()
	// zRouter.StrictSlash(true)
	zRouter.HandleFunc("/", zInfo).Methods("GET")
	zRouter.HandleFunc("/info", zInfo).Methods("GET")                                                        //1
	zRouter.HandleFunc("/task/", zDB.zTaskCreate).Methods("POST")                                            //2
	zRouter.HandleFunc("/task/", zDB.zTaskGetAll).Methods("GET")                                             //3
	zRouter.HandleFunc("/task/{id:[0-9]+}", zDB.zTaskDelete).Methods("DELETE")                               //4
	zRouter.HandleFunc("/task", zDB.zTaskDeleteAll).Methods("DELETE")                                        //5
	zRouter.HandleFunc("/tag/{name}", zDB.zTaskGetByTag).Methods("GET")                                      //6
	zRouter.HandleFunc("/due/{yyyy:[0-9]{4}}/{mm:[0-9]{2}}/{dd:[0-9]{2}}", zDB.zTaskGetByYMD).Methods("GET") //7
	zRouter.HandleFunc("/due/{yyyy:[0-9]{4}}/{mm:[0-9]{2}}", zDB.zTaskGetByYM).Methods("GET")                //8
	zRouter.HandleFunc("/due/{yyyy:[0-9]{4}}", zDB.zTaskGetByY).Methods("GET")                               //9

	zLog("Started at " + zHostPort + "\n")
	log.Fatal(http.ListenAndServe(zHostPort, zRouter)) // Запускаем web-сервер в режиме "слушания"

}

func zInfo(w http.ResponseWriter, r *http.Request) {
	zLog("zInfo")
	fmt.Fprintf(w, "\n<h1>TaskStoreAPI</h1>")
	fmt.Fprintf(w, "\n<p>1) <a href='/info'   target='_blank'>/info</a>       // Информация об API</p>")
	fmt.Fprintf(w, `  <p>2) /task/ // Сохранить новую задачу (запустить в консоли) // curl --header "Content-Type: application/json" --request POST --data '{"text":"test_text","tags":["tag1","tag2"],"due":"2024-03-04T20:21:22.123Z"}' http://127.0.0.1:5555/task/  </p>`)
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
