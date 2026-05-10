package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Book struct {
	ID     int
	Title  string
	Author string
}

var books []Book
var nextID = 1

func GetAllBooks(writer http.ResponseWriter, request *http.Request) {
	if len(books) == 0 {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode("Book Not found")
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(books)
}

func GetBookByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idBook := params.ByName("id")
	id, err := strconv.Atoi(idBook)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode("Book Not found")
		return
	}

	for _, book := range books {
		if book.ID == id {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(book)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)
	json.NewEncoder(writer).Encode("Book Not found")
}

func CreateBook(writer http.ResponseWriter, request *http.Request) {
	var book Book

	err := json.NewDecoder(request.Body).Decode(&book)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": "body kosong"})
		return
	}

	if book.Title == "" || book.Author == "" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": "title dan author wajib diisi"})
		return
	}

	book.ID = nextID
	nextID++
	books = append(books, book)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(book)
}

func DeleteBookByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idBook := params.ByName("id")
	id, err := strconv.Atoi(idBook)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode("Book Not found")
		return
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode("success delete book")
			return
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)
	json.NewEncoder(writer).Encode("Book Not found")
}
