package handler

import (
	"github.com/dorakueyon/gochain/blockchain"
	"net/http"
)

type Env struct {
	bc *blockchain.Blockchain
}

type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

	}

}
