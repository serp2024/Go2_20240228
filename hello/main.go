package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello, world!")

	client := &http.Client{Timeout: time.Minute}
	resOne, err := client.Get("http://golang.org/doc")
	if err != nil {
		log.Fatal(err)
	}

	// Свойство
	fmt.Println(resOne.Status, resOne.StatusCode, resOne.Request.URL)
	fmt.Println("-----------")
	fmt.Println(resOne.Request.Header)
	fmt.Println("-----------")
	fmt.Println(resOne.Request.Response.Header)
	fmt.Println("-----------")
	body, _ := io.ReadAll(resOne.Body)
	resOne.Body.Close()
	file, err := os.Create("out_site.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		log.Fatal(err)
	}

}
