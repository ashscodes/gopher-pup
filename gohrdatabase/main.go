package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	ConnectDatabase()
	router := NewRouter()
	fmt.Println("Listening on http://localhost:12345")
	fmt.Println("Press 'CTRL + C' to stop server.")
	log.Fatal(http.ListenAndServe(":12345", router))
}