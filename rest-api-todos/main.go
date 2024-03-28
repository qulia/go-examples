package main

import (
	"log"
	"net/http"

	"github.com/qulia/go-examples/rest-api-todos/middleware"
	"github.com/qulia/go-qulia/lib/set"
)

var db set.SetFlex[Todo, int64]

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", errorHandler(getAllTodos))
	mux.HandleFunc("DELETE /todos", errorHandler(deleteTodo))
	mux.HandleFunc("POST /todos", errorHandler(createTodo))
	mux.HandleFunc("GET /todos/{id}", errorHandler(getTodo))
	handler := middleware.Logger(mux)

	s := &http.Server{
		Addr:    ":3000",
		Handler: handler,
	}
	log.Fatal(s.ListenAndServe())
}

func errorHandler(f func(w http.ResponseWriter, r *http.Request) (int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := f(w, r)
		if err != nil {
			log.Printf("could not process request %v", err.Error())
			http.Error(w, err.Error(), code)
		}
	}
}
