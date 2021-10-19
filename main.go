package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/articles", getArticles)
	router.HandleFunc("/article", createArticle).Methods("POST")
	router.HandleFunc("/article/{id:[0-9]+}", getArticle).Methods("GET")
	router.HandleFunc("/article/{id:[0-9]+}", deleteArticle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
