package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

var books []Book
var nextID = 1

func GetAllBooks(writer http.ResponseWriter, request *http.Request) {
	if len(books) == 0 {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		encoder := json.NewEncoder(writer)
		_ = encoder.Encode(&Response{
			Status:  "error",
			Message: "No books found",
		})
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	_ = encoder.Encode(&Response{
		Status:  "success",
		Message: "All books found",
		Data:    books,
	})
}

func GetBookByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idBook := params.ByName("id")
	id, err := strconv.Atoi(idBook)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(writer)
		_ = encoder.Encode(&Response{
			Status:  "error",
			Message: "Invalid Book ID",
		})
		return
	}

	for _, book := range books {
		if book.ID == id {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			encoder := json.NewEncoder(writer)
			_ = encoder.Encode(&Response{
				Status:  "success",
				Message: "Book found",
				Data:    book,
			})
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)
	encoder := json.NewEncoder(writer)
	_ = encoder.Encode(&Response{
		Status:  "error",
		Message: "No books found",
	})
}

func CreateBook(writer http.ResponseWriter, request *http.Request) {
	var book Book

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&book)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(writer)
		_ = encoder.Encode(&Response{
			Status:  "error",
			Message: "Invalid format body",
		})
		return
	}

	if book.Title == "" || book.Author == "" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(writer)
		_ = encoder.Encode(&Response{
			Status:  "error",
			Message: "title and author cannot be empty",
		})
		return
	}

	book.ID = nextID
	nextID++
	books = append(books, book)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(writer)
	_ = encoder.Encode(&Response{
		Status:  "success",
		Message: "Created book",
		Data:    books,
	})
}

func DeleteBookByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idBook := params.ByName("id")
	id, err := strconv.Atoi(idBook)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(writer)
		_ = encoder.Encode(&Response{
			Status:  "error",
			Message: "Invalid Book ID",
		})
		return
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			encoder := json.NewEncoder(writer)
			_ = encoder.Encode(&Response{
				Status:  "success",
				Message: "Book deleted",
			})
			return
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)
	encoder := json.NewEncoder(writer)
	_ = encoder.Encode(&Response{
		Status:  "error",
		Message: "No books found",
	})
}
