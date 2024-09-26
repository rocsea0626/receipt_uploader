package main

import (
	"fmt"
	"net/http"
	"receipt_uploader/constants"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/health", helloHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(constants.PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
