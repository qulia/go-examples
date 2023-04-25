package main

type Todo struct {
	Name        string `json:"name"`
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (td Todo) Key() int64 {
	return td.Id
}
