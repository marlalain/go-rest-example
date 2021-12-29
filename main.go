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

func getBooks(w http.ResponseWriter, _ *http.Request) {
	log.Println("Returning list of books...")
	w.Header().Set("Content-Type", "application/json")
	sort.Sort(Books(books))
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			log.Println("Returning single book...")
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
	}

	err := json.NewEncoder(w).Encode(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	lastBookID, err := strconv.Atoi(books[len(books)-1].ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Creating book...")
	book.ID = strconv.Itoa(lastBookID + 1)
	books = append(books, book)
	encodeErr := json.NewEncoder(w).Encode(book)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			log.Println("Updating book...")
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = item.ID
			books = append(books, book)
			err := json.NewEncoder(w).Encode(book)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			return
		}
	}

	err := json.NewEncoder(w).Encode(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			log.Println("Deleting book...")
			books = append(books[:index], books[index+1:]...)
			break
		}
	}

	err := json.NewEncoder(w).Encode(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func mockData(data []Book) []Book {
	log.Println("Creating mock data...")

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
