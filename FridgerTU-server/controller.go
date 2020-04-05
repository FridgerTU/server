package main

import (
	"fmt"
	"net/http"
	"os"
)

func NewController() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basePathHandler)
	mux.HandleFunc("/api/v1/recipes", recipesHandler)
	mux.HandleFunc("/api/v1/recipe", recipeHandler)
	mux.HandleFunc("/api/v1/random", randomHandler)
	return mux
}

func basePathHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte("hello")); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func recipesHandler(writer http.ResponseWriter, request *http.Request) {
	//TODO implement
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte("recipes")); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func recipeHandler(writer http.ResponseWriter, request *http.Request) {
	//TODO implement
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte("recipe")); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func randomHandler(writer http.ResponseWriter, request *http.Request) {
	url := "https://www.themealdb.com/api/json/v1/1/random.php"
	bytes, err := executeGetRequest(url)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write(bytes); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	//TODO print a custom json
}
