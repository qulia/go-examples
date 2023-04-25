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
func getTodo(writer http.ResponseWriter, request *http.Request) (int, error) {
	if ids, ok := mux.Vars(request)["id"]; ok && len(ids) != 0 {
		id, err := strconv.Atoi(ids)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		td := db.GetWithKey(int64(id))
		if td != *new(Todo) {
			err = json.NewEncoder(writer).Encode(td)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			return http.StatusNotFound, fmt.Errorf("not found")
		}
	}
	return http.StatusOK, nil
}

/*
curl localhost:3000/todos
*/
func getAllTodos(w http.ResponseWriter, _ *http.Request) (int, error) {
	err := json.NewEncoder(w).Encode(db.ToSlice())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

/*
	curl \
	    --header "Content-type: application/json" \
	    --request POST \
	    --data '{"name":"todo","id":2,"title":"t2","description":"d2"}' \
	    http://localhost:3000/todo
*/
func createTodo(w http.ResponseWriter, r *http.Request) (int, error) {
	var td Todo
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&td)
	if err != nil {
		return http.StatusBadRequest, err
	}
	db.Add(td)
	err = json.NewEncoder(w).Encode(td)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

/*
	curl \
	    --header "Content-type: application/json" \
	    --request DELETE \
	    --data '{"name":"todo","id":2}' \
	    http://localhost:3000/todo
*/
func deleteTodo(_ http.ResponseWriter, r *http.Request) (int, error) {
	var td Todo
	err := json.NewDecoder(r.Body).Decode(&td)
	if err != nil {
		return http.StatusBadRequest, err
	}
	if db.Contains(td) {
		db.Remove(td)
	} else {
		return http.StatusNotFound, fmt.Errorf("not found")
	}

	return http.StatusOK, nil
}
