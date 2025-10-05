package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	pr, pw := io.Pipe()
	defer func() {
		pr.Close()
		pw.Close()
	}()

	var writeErr error
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				log.Print("ticker.C")
				_, err := pw.Write([]byte("hello\n"))
				if err != nil {
					writeErr = fmt.Errorf("unable to write data to server: %s", err)
					stop()
					return
				}
			case <-ctx.Done():
				log.Print("ctx.Done()")
				return
			}
		}
	}()

	//reqCtx, reqCtxCancel := context.WithTimeout(ctx, 2*time.Second)
	//defer reqCtxCancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://localhost:8080/api/full-duplex", pr)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2: true,
	}
	client := http.Client{
		Transport: tr,
	}

	//ch := make(chan *http.Response, 1)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to request: %s", err)
	}
	log.Println("sent request")

	var readErr error
	go func() {
		bufReader := bufio.NewReader(res.Body)
		for {
			select {
			case <-ctx.Done():
				log.Print("ctx.Done() 2")
				return
			default:
				log.Print("read")
				req, err := bufReader.ReadString('\n')
				if err != nil {
					readErr = fmt.Errorf("unable to ready response body: %s", err)
					stop()
					return
				}
				log.Print(req)
			}
		}
	}()

	<-ctx.Done()
	return errors.Join(readErr, writeErr)
}
