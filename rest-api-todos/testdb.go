package main

import (
	"github.com/qulia/go-qulia/lib/set"
)

func init() {
	db = set.NewSetFlex[Todo, int64]()

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
