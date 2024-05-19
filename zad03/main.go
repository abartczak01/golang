package main

import (
	"log"
	"net/http"
	"zad03/handlers"
)

func main() {
    http.HandleFunc("/posts", handlers.PostsHandler)
    http.HandleFunc("/posts/", handlers.PostHandler)

    log.Println("Server is running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}