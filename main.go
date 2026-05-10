package main

import (
	"api-buku/handler"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Middleware struct {
	http.Handler
}

func (middleware Middleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	middleware.Handler.ServeHTTP(writer, request)
	fmt.Println(request.Method, request.URL.Path)
}

func main() {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Not found")
	})

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, p interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server error")
	}

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Method not allowed")
	})

	router.GET("/books/", func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		handler.GetAllBooks(writer, request)
	})
	router.GET("/books/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handler.GetBookByID(writer, request, params)
	})
	router.POST("/books", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handler.CreateBook(writer, request)
	})
	router.DELETE("/books/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handler.DeleteBookByID(writer, request, params)
	})

	logMiddleware := &Middleware{Handler: router}

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: logMiddleware,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
