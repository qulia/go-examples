package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qulia/go-qulia/lib/set"
)

var db *set.SliceSet

func main() {
	gRouter := mux.NewRouter().StrictSlash(true)
	gRouter.HandleFunc("/todos", errorHandler(getAllTodos))
	gRouter.HandleFunc("/todo", errorHandler(deleteTodo)).Methods("DELETE")
	gRouter.HandleFunc("/todo", errorHandler(createTodo)).Methods("POST")
	gRouter.HandleFunc("/todos/{id}", errorHandler(getTodo))
	gRouter.Use(loggingMiddleware)
	s := &http.Server{
		Addr:    ":3000",
		Handler: gRouter,
	}
	log.Fatal(s.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s:%s\n", r.RequestURI, r.Method)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func errorHandler(f func(w http.ResponseWriter, r *http.Request) (error, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err, code := f(w, r)
		if err != nil {
			log.Printf("could not process request %v", err.Error())
			if code > 0 {
				http.Error(w, err.Error(), code)
			} else {
				// Do not pass internal details to client
				http.Error(w, "Error occurred", http.StatusInternalServerError)
			}
		}
	}
}
