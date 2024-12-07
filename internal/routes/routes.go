package routes

import (
	"fmt"
	"go-test/internal/generate"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func GenerateSQL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp, err := generate.GenerateHandler(w, r, ps)
	if err != nil {
		// TODO: handle error
	}

	for _, d := range resp {
		fmt.Fprintf(w, "%v", d)
	}
}
