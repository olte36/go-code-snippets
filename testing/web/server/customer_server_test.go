package server_test

import (
	"encoding/json"
	"go-code-patterns/testing/web/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerServer(t *testing.T) {
	expected := []server.Customer{
		{
			Name:  "Alice",
			Email: "alice@mail.com",
		},
		{
			Name:  "Bob",
			Email: "bob@mail.com",
		},
	}

	request := httptest.NewRequest(http.MethodGet, "/customers", nil)
	recorder := httptest.NewRecorder()

	srv := server.NewCustomerServer(8080)
	srv.Handler.ServeHTTP(recorder, request)

	var customers []server.Customer
	err := json.NewDecoder(recorder.Body).Decode(&customers)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	assert.ElementsMatch(t, expected, customers)
}
