package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	fmt.Println("Hello welcome to log service!")
	http.HandleFunc("/log", printLog)
	http.ListenAndServe(":5003", nil)
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// d1 := []byte("Hello GET request\n")
	// err := ioutil.WriteFile("logfile", d1, 0644)

	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Receive a GET request")

	fmt.Fprintf(w, "Hello")
}

func printLog(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}

	log.Println(string(body))

	d1 := []byte(string(body))
	err = ioutil.WriteFile("logfile", d1, 0644)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "ok")
}
