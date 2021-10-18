package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func getArticles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Articles)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}
	a := Article{Id: id}
	if err := a.deleteArticle(); err != nil {
		switch err {
		case errors.New("not found"):
			respondWithError(w, http.StatusNotFound, "Article not found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Article not found")
		}
	}
	respondWithJSON(w, http.StatusAccepted, map[string]string{"result": "success"})
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	var a Article
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&a); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := a.createArticle(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, a)
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}
	a := Article{Id: id}
	if err := a.getArticle(); err != nil {
		switch err {
		case errors.New("not found"):
			respondWithError(w, http.StatusNotFound, "Article not found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Article not found")
		}
	}
	respondWithJSON(w, http.StatusOK, a)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type Article struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []Article

func (a *Article) getArticle() error {
	for _, article := range Articles {
		if article.Id == a.Id {
			a.Desc = article.Desc
			a.Title = article.Title
			a.Content = article.Content
		}
	}
	return errors.New("not found")
}

// func (a *Article) getArticle() error {
// 	for _, article := range Articles {
// 		if article.Id == a.Id {
// 			a.Desc = article.Desc
// 			a.Title = article.Title
// 			a.Content = article.Content
// 		}
// 	}
// 	return errors.New("not found")
// }

func (a *Article) createArticle() error {
	Articles = append(Articles, *a)
	return nil
}

func (a *Article) deleteArticle() error {
	for i, article := range Articles {
		if article.Id == a.Id {
			Articles = append(Articles[:i], Articles[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []Article{
		{Id: 1, Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		{Id: 2, Title: "Hello2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}
