package client_test

import (
	"fmt"
	"go-code-patterns/testing/web/client"
)

// https://pkg.go.dev/testing#hdr-Examples
func ExampleGitHubClient_GetTags() {
	c := client.New()
	tags, err := c.GetTags("golang", "go")
	if err != nil {
		fmt.Println(err)
		return
	}
	var go123 client.Tag
	for _, t := range tags {
		if t.Ref == "refs/tags/go1.23.0" {
			go123 = t
			break
		}
	}
	fmt.Printf("%+v", go123)
	// Output:
	// {Ref:refs/tags/go1.23.0 NodeId:MDM6UmVmMjMwOTY5NTk6cmVmcy90YWdzL2dvMS4yMy4w Url:https://api.github.com/repos/golang/go/git/refs/tags/go1.23.0 Object:{Sha:6885bad7dd86880be6929c02085e5c7a67ff2887 Type:commit Url:https://api.github.com/repos/golang/go/git/commits/6885bad7dd86880be6929c02085e5c7a67ff2887}}
}
