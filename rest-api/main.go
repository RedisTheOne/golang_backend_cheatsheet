package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//Book Struct
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

//Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//GetBooks Function
func getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

//Get Single Book Function
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

//Create Book Function
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = uuid.New().String()
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

//Update Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			var book Book
			books = append(books[:index], books[index+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

//Delete Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

//Using Headers Example
func authorized(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	authArr := strings.Split(auth, "Bearer ")
	fmt.Println(authArr[1])
}

//Init Books Variable
var books []Book

func main() {
	r := mux.NewRouter()
	books = append(books, Book{ID: "1", Isbn: "123456789", Title: "First Book", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "123456789", Title: "Second Book", Author: &Author{Firstname: "John", Lastname: "Doe"}})

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	r.HandleFunc("/api/auth", authorized).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
