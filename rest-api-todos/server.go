package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
 curl localhost:3000/todos/1
*/
func getTodo(writer http.ResponseWriter, request *http.Request) (error, int) {
	if ids, ok := mux.Vars(request)["id"]; ok && len(ids) != 0 {
		id, err := strconv.Atoi(ids)
		if err != nil {
			return err, -1
		}

		td := db.GetItemForKey(Todo{Id: int64(id)})
		if td != nil {
			err = json.NewEncoder(writer).Encode(td)
			if err != nil {
				return err, -1
			}
		} else {
			return fmt.Errorf("not found"), http.StatusNotFound
		}
	}
	return nil, 0
}

/*
 curl localhost:3000/todos
*/
func getAllTodos(w http.ResponseWriter, _ *http.Request) (error, int) {
	return json.NewEncoder(w).Encode(db.GetSlice()), -1
}

/*
curl \
    --header "Content-type: application/json" \
    --request POST \
    --data '{"name":"todo","id":2,"title":"t2","description":"d2"}' \
    http://localhost:3000/todo
*/
func createTodo(w http.ResponseWriter, r *http.Request) (error, int) {
	var td Todo
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&td)
	if err != nil {
		return err, http.StatusBadRequest
	}
	db.Add(td)
	err = json.NewEncoder(w).Encode(td)
	if err != nil {
		return err, -1
	}

	return nil, 0
}

/*
curl \
    --header "Content-type: application/json" \
    --request DELETE \
    --data '{"name":"todo","id":2}' \
    http://localhost:3000/todo
*/
func deleteTodo(_ http.ResponseWriter, r *http.Request) (error, int) {
	var td Todo
	err := json.NewDecoder(r.Body).Decode(&td)
	if err != nil {
		return err, http.StatusBadRequest
	}
	if db.ContainsKeyFor(td) {
		db.Remove(td)
	} else {
		return fmt.Errorf("not found"), http.StatusNotFound
	}

	return nil, 0
}
