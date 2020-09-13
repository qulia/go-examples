package main

import (
	"fmt"

	"github.com/qulia/go-qulia/lib/set"
)

func init() {
	db = set.NewSliceSet(func(it interface{}) string {
		return fmt.Sprintf("%d", it.(Todo).Id)
	})

	db.Add(Todo{
		Name:        "todo",
		Id:          0,
		Title:       "t0",
		Description: "d0",
	})
	db.Add(Todo{
		Name:        "todo",
		Id:          1,
		Title:       "t1",
		Description: "d1",
	})
}
