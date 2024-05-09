package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/feliux/go-bdd/http/bookstore"
)

type App struct {
	store *bookstore.Store
}

func (a *App) getBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	if r.Method == http.MethodGet {
		book, err := a.store.GetBook(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		bookJson, err := json.Marshal(book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		_, _ = w.Write(bookJson)
	}
}

func (a *App) addBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var book bookstore.Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = a.store.AddBook(&book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func main() {
	app := App{store: &bookstore.Store{}}
	// curl -X POST localhost:8080/books -d '{"id":"1", "title":"El Principito", "author":"Antoine de Saint-Exupery"}'
	http.HandleFunc("/books", app.addBook)
	// curl -X GET localhost:8080/books/1
	http.HandleFunc("/books/", app.getBook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
