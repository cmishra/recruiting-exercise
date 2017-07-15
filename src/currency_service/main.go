package main

import (
	"net/http"
	"log"
	"io"
)


func requestHandler(w http.ResponseWriter, r *http.Request) {
//	log.Println(r.URL.RawQuery)
	io.WriteString(w, r.URL.RawQuery + "\n")
}

const (
	port = ":9000"
)


func main() {
	log.Println("Listening to " + port)
	http.HandleFunc("/rates", requestHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}
