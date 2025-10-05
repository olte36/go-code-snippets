package client_test

import (
	"encoding/json"
	"go-code-patterns/testing/web/client"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	expectedResp := []client.Tag{
		{
			Ref:    "refs/tags/go1.0.1",
			NodeId: "MDM6UmVmMjMwOTY5NTk6cmVmcy90YWdzL2dvMS4wLjE=",
			Url:    "https://api.github.com/repos/golang/go/git/refs/tags/go1.0.1",
			Object: client.Object{
				Sha:  "2fffba7fe19690e038314d17a117d6b87979c89f",
				Type: "commit",
				Url:  "https://api.github.com/repos/golang/go/git/commits/2fffba7fe19690e038314d17a117d6b87979c89f",
			},
		},
		{
			Ref:    "refs/tags/go1.0.2",
			NodeId: "MDM6UmVmMjMwOTY5NTk6cmVmcy90YWdzL2dvMS4wLjI=",
			Url:    "https://api.github.com/repos/golang/go/git/refs/tags/go1.0.2",
			Object: client.Object{
				Sha:  "cb6c6570b73a1c4d19cad94570ed277f7dae55ac",
				Type: "commit",
				Url:  "https://api.github.com/repos/golang/go/git/commits/cb6c6570b73a1c4d19cad94570ed277f7dae55ac",
			},
		},
	}
	rawResp, err := json.Marshal(expectedResp)
	if !assert.Nil(t, err) {
		t.Fatal("unable to marshal")
	}
	srv := mockServer(rawResp)
	defer srv.Close()

	ghClient := client.NewWithBaseUrl(srv.URL)
	tags, err := ghClient.GetTags("golang", "go")
	if !assert.Nil(t, err) {
		t.Fatal("unable to get tags")
	}

	assert.ElementsMatch(t, expectedResp, tags)
}

func mockServer(respBody []byte) *httptest.Server {
	handlerF := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write(respBody)
		if err != nil {
			panic(err)
		}
	}
	return httptest.NewServer(http.HandlerFunc(handlerF))
}
