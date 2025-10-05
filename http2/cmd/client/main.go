package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	//reqCtx, reqCtxCancel := context.WithTimeout(ctx, 2*time.Second)
	//defer reqCtxCancel()

	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2: true,
	}
	client := http.Client{
		Transport: tr,
	}

	endpoint := "https://localhost:8080/api/full-duplex"

	buffer := bytes.NewBufferString("2\n3\n4\n5\n")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buffer)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to request: %s", err)
	}
	log.Println("sent request")

	bufReader := bufio.NewReader(res.Body)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			res, err := bufReader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
			log.Print(res)
		}
	}
}
