package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error

// func (h *Handler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	err := h(w, r, ps)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError) // TOOD: this should be different status codes, can the func return a code
// 	}
// }
