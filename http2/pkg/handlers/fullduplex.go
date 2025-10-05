package handlers

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"go-code-patterns/http2/pkg/respctrl"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var _ http.Handler = fullDuplexHandler{}
var _ error = (*httpErr)(nil)

type fullDuplexHandler struct {
	fullDuplexForHttp1 bool
	flushResp          bool
}

type httpErr struct {
	status int
	err    error
}

type loggedReader struct {
	r io.Reader
}

func NewFullDuplexHandler(flushResp, fullDuplexForHttp1 bool) http.Handler {
	return fullDuplexHandler{
		flushResp:          flushResp,
		fullDuplexForHttp1: fullDuplexForHttp1,
	}
}

func (f fullDuplexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var wg sync.WaitGroup
	numsCh := make(chan int)
	errCh := make(chan error, 1)

	const numWorkers = 3
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			if err := f.sendResponse(ctx, w, numsCh); err != nil {
				select {
				case errCh <- err:
				default:
					slog.Error("error while sending response", slog.Any("err", err))
				}
			}
			cancel()
		}()
	}

	readBodyConcurently := true
	if r.ProtoMajor == 1 {
		if f.fullDuplexForHttp1 {
			rc := respctrl.NewImprovedResponseController(w)
			if err := rc.EnableFullDuplex(); err != nil {
				readBodyConcurently = false
				slog.Warn("error while enabling full duplex", slog.Any("err", err))
			}
		} else {
			readBodyConcurently = false
		}
	}

	if readBodyConcurently {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := readBody(ctx, r, numsCh, false); err != nil {
				select {
				case errCh <- err:
				default:
					slog.Error("error while reading body concurrently", slog.Any("err", err))
				}
			}
			cancel()
		}()
	} else {
		if err := readBody(ctx, r, numsCh, true); err != nil {
			select {
			case errCh <- err:
			default:
				slog.Error("error while reading body", slog.Any("err", err))
			}
		}
		cancel()
	}

	wg.Wait()

	select {
	case err := <-errCh:
		status := http.StatusInternalServerError
		if httpErr, ok := err.(*httpErr); ok {
			status = httpErr.status
		}
		http.Error(w, err.Error(), status)
	default:
	}
}

func (f fullDuplexHandler) sendResponse(ctx context.Context, w http.ResponseWriter, numsCh chan int) error {
	rc := respctrl.NewImprovedResponseController(w)
	flushSupported := f.flushResp
	for {
		select {
		case <-ctx.Done():
			return nil
		case num := <-numsCh:
			// simulate some work
			mill := rand.IntN(3000) + 1000
			time.Sleep(time.Duration(mill) * time.Millisecond)

			fmt.Fprintf(w, "%d^2=%d\n", num, num*num)

			if flushSupported {
				if err := rc.Flush(); err != nil {
					if errors.Is(err, http.ErrNotSupported) {
						slog.Warn("Flushing is not supported")
						flushSupported = false
					} else {
						return &httpErr{
							status: http.StatusInternalServerError,
							err:    err,
						}
					}
				}
			}
		}
	}
}

func readBody(ctx context.Context, r *http.Request, numsCh chan int, readAll bool) error {
	var reader io.Reader = loggedReader{r: r.Body}
	if readAll {
		body, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("unable to read the whole body: %w", err)
		}
		reader = bytes.NewBuffer(body)
	}
	bufReader := bufio.NewReaderSize(reader, 16)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			strNum, err := bufReader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return &httpErr{
					status: http.StatusInternalServerError,
					err:    err,
				}
			}

			// simulate some work
			//time.Sleep(time.Second)

			num, err := strconv.Atoi(strNum[:len(strNum)-1])
			if err != nil {
				return &httpErr{
					status: http.StatusBadRequest,
					err:    err,
				}
			}
			// prevent goroutine leak
			timer := time.NewTimer(15 * time.Second)
			select {
			case numsCh <- num:
			case <-timer.C:
				return errors.New("stop reading the body, processing takes too long")
			}
		}
	}
}

func (h *httpErr) Error() string {
	return fmt.Sprintf("%d %s", h.status, h.err)
}

func (l loggedReader) Read(p []byte) (int, error) {
	n, err := l.r.Read(p)
	slog.Info("Read body to buf", slog.Int("bytes", n))
	return n, err
}
