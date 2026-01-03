package main

import (
	"log"
	"net/http"
	"os"
)

func startServer() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/health", handleHealthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7878"
	}

	addr := "0.0.0.0:" + port
	log.Println("Listening on ", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
