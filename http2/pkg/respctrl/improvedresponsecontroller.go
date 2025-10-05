package respctrl

import (
	"fmt"
	"net/http"
)

type ImprovedResponseController struct {
	*http.ResponseController
	rw http.ResponseWriter
}

func NewImprovedResponseController(rw http.ResponseWriter) *ImprovedResponseController {
	return &ImprovedResponseController{
		ResponseController: http.NewResponseController(rw),
		rw:                 rw,
	}
}

func (i *ImprovedResponseController) Push(target string, opts *http.PushOptions) error {
	rw := i.rw
	for {
		switch t := rw.(type) {
		case http.Pusher:
			return t.Push(target, opts)
		case interface{ Unwrap() http.ResponseWriter }:
			rw = t.Unwrap()
		default:
			return fmt.Errorf("%w", http.ErrNotSupported)
		}
	}
}
