package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/signal"
	"strconv"
	"syscall"
)

const delim = '\n'

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var serverAddr string
	flag.StringVar(&serverAddr, "addr", "https://localhost:8080", "")
	flag.Parse()

	// prepate endpoint
	serverUrl, err := url.Parse(serverAddr)
	if err != nil {
		return fmt.Errorf("invalid server address %s: %w", serverAddr, err)
	}
	serverUrl.Path = "/api/full-duplex"

	// use HTTP2 by default, skip certificate checking
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2: true,
	}
	client := http.Client{
		Transport: tr,
	}

	// prepare request body
	var buffer bytes.Buffer
	for num := range 16 {
		buffer.WriteString(strconv.Itoa(num))
		buffer.WriteRune(delim)
	}

	// preapre and send request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverUrl.String(), &buffer)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to request: %s", err)
	}

	// read response body
	bufReader := bufio.NewReader(res.Body)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			res, err := bufReader.ReadString(delim)
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
