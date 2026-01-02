package main

import "net/http"

func startServer() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/health", handleHealthCheck)

	http.ListenAndServe(":7878", nil)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
