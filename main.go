package main

import (
	"api-buku/handler"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Middleware struct {
	http.Handler
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (logResponseWriter *LoggingResponseWriter) WriteHeader(code int) {
	logResponseWriter.statusCode = code
	logResponseWriter.ResponseWriter.WriteHeader(code)
}

func (middleware Middleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logResponseWriter := &LoggingResponseWriter{
		ResponseWriter: writer,
		statusCode:     http.StatusOK,
	}

	start := time.Now()

	middleware.Handler.ServeHTTP(logResponseWriter, request)
	fmt.Println(
		request.Method,
		request.URL.Path,
		logResponseWriter.statusCode,
		time.Since(start),
	)
}

func main() {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		encoder := json.NewEncoder(w)
		_ = encoder.Encode(&handler.Response{
			Status:  "error",
			Message: "not found",
		})
	})

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, p interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(w)
		_ = encoder.Encode(&handler.Response{
			Status:  "error",
			Message: "internal server error",
		})
	}

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder := json.NewEncoder(w)
		_ = encoder.Encode(&handler.Response{
			Status:  "error",
			Message: "method not allowed",
		})
	})

	router.GET("/books", func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
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
