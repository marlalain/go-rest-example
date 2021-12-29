package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type Book struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Books []Book

type Author struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var books []Book

func (books Books) Len() int {
	return len(books)
}

func (books Books) Less(i, j int) bool {
	return books[i].ID < books[j].ID
}

func (books Books) Swap(i, j int) {
	books[i], books[j] = books[j], books[i]
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sort.Sort(Books(books))
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(nil)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	lastBookID, err := strconv.Atoi(books[len(books)-1].ID)
	if err != nil {
		return
	}
	book.ID = strconv.Itoa(lastBookID + 1)
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = item.ID
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	json.NewEncoder(w).Encode(nil)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(nil)
}

func mockData(data []Book) []Book {
	data = append(data, Book{
		ID:    "1",
		ISBN:  "1111",
		Title: "My Book",
		Author: &Author{
			FirstName: "John",
			LastName:  "Doe",
		},
	})

	data = append(data, Book{
		ID:    "2",
		ISBN:  "2222",
		Title: "The Better Book",
		Author: &Author{
			FirstName: "Will",
			LastName:  "Smith",
		},
	})

	return data
}

func main() {
	router := mux.NewRouter()
	log.Println("Starting server...")

	books = mockData(books)

	router.HandleFunc("/api/v1/books", getBooks).Methods("GET")
	router.HandleFunc("/api/v1/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/v1/books", createBook).Methods("POST")
	router.HandleFunc("/api/v1/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/v1/books/{id}", deleteBook).Methods("DELETE")

	log.Fatalln(http.ListenAndServe(":8000", router))
}
