package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	response string
	port     int
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, response)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)
	log.Printf("Server is running (port: %d)..", port)
	log.Printf("Response body for upcoming requests will be: %s", response)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func init() {
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.StringVar(&response, "response", "", "String to be sent in response body")
}
