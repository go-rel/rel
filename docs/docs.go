package main

// User model example.
type User struct {
	ID   int
	Name string
}

// Transaction model example.
type Transaction struct {
	ID     int
	Book   Book
	BookID int
	Status string
}
