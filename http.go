package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getArticles(w http.ResponseWriter, r *http.Request) {
	a, err := getArticlesFromRepo()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid article query")
	}
	respondWithJSON(w, http.StatusOK, a)
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
		a = Article{}
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
