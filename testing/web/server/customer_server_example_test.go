package server_test

import (
	"context"
	"errors"
	"fmt"
	"go-code-patterns/testing/web/server"
	"net/http"
	"os/signal"
	"syscall"
)

func ExampleNewCustomerServer() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv := server.NewCustomerServer(8080)
	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)

	//// Output:
	//
}
