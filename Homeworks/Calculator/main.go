/*
## Задача № 1
Написать API для указанных маршрутов(endpoints)
"/info"   // Информация об API
"/first"  // Случайное число
"/second" // Случайное число
"/add"    // Сумма двух случайных чисел
"/sub"    // Разность
"/mul"    // Произведение
"/div"    // Деление

*результат вернуть в виде JSON

"math/rand"
number := rand.Intn(100)
! не забудьте про Seed()


GET http://127.0.0.1:1234/first

GET http://127.0.0.1:1234/second

GET http://127.0.0.1:1234/add
GET http://127.0.0.1:1234/sub
GET http://127.0.0.1:1234/mul
GET http://127.0.0.1:1234/div
GET http://127.0.0.1:1234/info
*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ztCalc struct {
	First     int    `json:"first"`
	Second    int    `json:"second"`
	Operation string `json:"operation"`
	Result    string `json:"result"`
	ResultTxt string `json:"result_txt"`
}

func main() {
	zLog("help run:   SERVERPORT=8888 go run . \n")
	zPort := os.Getenv("SERVERPORT")
	if zPort == "" {
		zPort = "1234"
	}
	zHostPort := "0.0.0.0:" + zPort
	http.HandleFunc("/favicon.ico", zfavicon) // favicon for test from  chrome
	http.HandleFunc("/", zHome)               // Информация об API
	http.HandleFunc("/info", zInfo)           // Информация об API
	http.HandleFunc("/first", zFirst)         // Случайное число
	http.HandleFunc("/second", zSecond)       // Случайное число
	http.HandleFunc("/add", zRouteOpration)
	http.HandleFunc("/sub", zRouteOpration)
	http.HandleFunc("/mul", zRouteOpration)
	http.HandleFunc("/div", zRouteOpration)
	zLog("Started at " + zHostPort + "\n")
	log.Fatal(http.ListenAndServe(zHostPort, nil)) // Запускаем web-сервер в режиме "слушания"
}

// zCalcBegin для ответа set Content-Type=application/json  и  заполнение из прилетевшего от клиента json-a в переменную zReq
func zCalcBegin(w http.ResponseWriter, r *http.Request, zRCalc *ztCalc) {
	zLog("zCalcBegin ")
	if r.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&zRCalc) // заполнение из прилетевшего от клиента json-a в переменную  zRCalc (zReq)
		if err != nil {
			zErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
	}

	// cookie для работоспособности апи-шки в браузере
	calc2cookie_str, err := r.Cookie("calc2cookie")
	if err == nil {
		calc2cookie_bytes, _ := base64.StdEncoding.DecodeString(calc2cookie_str.Value)
		json.Unmarshal(calc2cookie_bytes, zRCalc)
	}

}

// zSendRes отправка ответа в body для API  и дублем в cookie для возможности пользоваться в браузере
func zSendRes(w http.ResponseWriter, zCalc ztCalc) {
	zLog("zSendRes ")
	calc2cookie_bytes, err := json.Marshal(zCalc)
	if err == nil {
		calc2cookie_str := base64.StdEncoding.EncodeToString(calc2cookie_bytes)
		cookie := http.Cookie{
			Name:     "calc2cookie",
			Value:    calc2cookie_str,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zCalc)

}

func zFirst(w http.ResponseWriter, r *http.Request) {
	zLog("zFirst ")
	var zCalc ztCalc
	zCalcBegin(w, r, &zCalc)
	zCalc.ResultTxt = "set first random. "
	zCalc.First = zGetRandomInt()
	zCalculate(&zCalc)
	zSendRes(w, zCalc)
	zLog(".\n")
}

func zSecond(w http.ResponseWriter, r *http.Request) {
	zLog("zSecond ")
	var zCalc ztCalc
	zCalcBegin(w, r, &zCalc)
	zCalc.ResultTxt = "set second random. "
	zCalc.Second = zGetRandomInt()
	zCalculate(&zCalc)
	zSendRes(w, zCalc)
	zLog(".\n")
}

func zRouteOpration(w http.ResponseWriter, r *http.Request) {
	zLog("zRouteOpration ")
	var zCalc ztCalc
	if r.Method == http.MethodGet {
		zCalcBegin(w, r, &zCalc)
		switch r.URL.Path {
		case "/add":
			zCalc.Operation = "add"
		case "/sub":
			zCalc.Operation = "sub"
		case "/mul":
			zCalc.Operation = "mul"
		case "/div":
			zCalc.Operation = "div"
		default:
			zCalc.Operation = ""
		}
		zCalc.ResultTxt = ""
		zCalculate(&zCalc)
		zSendRes(w, zCalc)
	} else {
		zErrorResponse(w, "Supported only GET http method", http.StatusBadRequest)
	}
	zLog(".\n")
}

func zCalculate(zCalc *ztCalc) {
	zLog("zCalculate ")
	switch zCalc.Operation {
	case "add":
		zCalc.Result = fmt.Sprint(float64(zCalc.First) + float64(zCalc.Second))
		zCalc.ResultTxt += " " + strconv.Itoa(zCalc.First) + " + " + strconv.Itoa(zCalc.Second) + " = " + fmt.Sprint(zCalc.Result) + " "
	case "sub":
		zCalc.Result = fmt.Sprint(float64(zCalc.First) - float64(zCalc.Second))
		zCalc.ResultTxt += " " + strconv.Itoa(zCalc.First) + " - " + strconv.Itoa(zCalc.Second) + " = " + fmt.Sprint(zCalc.Result) + " "
	case "mul":
		zCalc.Result = fmt.Sprint(float64(zCalc.First) * float64(zCalc.Second))
		zCalc.ResultTxt += " " + strconv.Itoa(zCalc.First) + " * " + strconv.Itoa(zCalc.Second) + " = " + fmt.Sprint(zCalc.Result) + " "
	case "div":
		if zCalc.Second != 0 {
			zCalc.Result = fmt.Sprint(float64(zCalc.First) / float64(zCalc.Second))
			zCalc.ResultTxt += " " + strconv.Itoa(zCalc.First) + " / " + strconv.Itoa(zCalc.Second) + " = " + fmt.Sprint(zCalc.Result) + " "
		} else {
			zCalc.Result = ""
			zCalc.ResultTxt += " " + strconv.Itoa(zCalc.First) + " / " + strconv.Itoa(zCalc.Second) + " ??? division by zero is not possible "
		}
	case "":
		zCalc.Result = ""
		zCalc.ResultTxt += " skip calc (operation not set) "
	default:
		zCalc.Result = ""
		zCalc.ResultTxt += " operation not supported "
	}
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

func zLog(txt string) {
	vCurrTime := time.Now()
	if txt != ".\n" {
		fmt.Print(" ; ", vCurrTime.Format("15:04:05.000"), " ")
	}
	fmt.Print(txt)
}

func zGetRandomIntRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func zGetRandomInt() int {
	zRes := zGetRandomIntRange(0, 100)
	zLog("zGetRandomInt " + strconv.Itoa(zRes))
	return zRes
}

func zHome(w http.ResponseWriter, r *http.Request) {
	zLog("zHome ")
	zInfo(w, r)
}

// zInfo   -- Информация об API
func zInfo(w http.ResponseWriter, r *http.Request) {
	zLog("zInfo")
	fmt.Fprintf(w, "\n<h1>Calculator API</h1>")
	fmt.Fprintf(w, "\n<p><a href='/info'   target='_blank'>/info</a>   // Информация об API</p>")
	fmt.Fprintf(w, "\n<p><a href='/first'  target='_blank'>/first</a>  // Случайное число</p>")
	fmt.Fprintf(w, "\n<p><a href='/second' target='_blank'>/second</a> // Случайное число</p>")
	fmt.Fprintf(w, "\n<p><a href='/add'    target='_blank'>/add</a>    // Сумма двух случайных чисел</p>")
	fmt.Fprintf(w, "\n<p><a href='/sub'    target='_blank'>/sub</a>    // Разность</p>")
	fmt.Fprintf(w, "\n<p><a href='/mul'    target='_blank'>/mul</a>    // Произведение</p>")
	fmt.Fprintf(w, "\n<p><a href='/div'    target='_blank'>/div</a>    // Деление</p>")
	zLog(".\n")
}

func zfavicon(w http.ResponseWriter, r *http.Request) {
	zLog("zfavicon")
	//w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "image/x-icon") //"image/vnd.microsoft.icon")
	zIco := "\x00\x00\x01\x00\x01\x00\x10\x10\x10\x00\x01\x00\x04\x00\x28\x01\x00\x00\x16\x00\x00\x00\x28\x00\x00\x00\x10\x00\x00\x00\x20\x00\x00\x00\x01\x00\x04\x00\x00\x00\x00\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x10\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00\xFF\xFF\x00\x00"
	w.Write([]byte(zIco))
	zLog(".\n")
}
