package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Customer struct {
	Name  string
	Email string
}

func NewCustomerServer(port uint) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/customers", http.HandlerFunc(listCustomers))
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return &srv
}

func listCustomers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	customers := []Customer{
		{
			Name:  "Alice",
			Email: "alice@mail.com",
		},
		{
			Name:  "Bob",
			Email: "bob@mail.com",
		},
	}
	err := json.NewEncoder(w).Encode(customers)
	if err != nil {
		w.WriteHeader(500)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			panic(err)
		}
	}
}
