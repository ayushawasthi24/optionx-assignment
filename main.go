package main

import (
	"log"
	"net/http"
	"optionx-assignment/server"
)

func main() {
	s := server.NewServer()
	http.HandleFunc("/ws", s.HandleConnections)
	log.Println("WebSockets server started on :8080")
	http.ListenAndServe(":8080", nil)
}
