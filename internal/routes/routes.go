package routes

import (
	"encoding/json"
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
	resp := generate.GenerateHandler(r)
	if resp.Status != http.StatusOK {
		http.Error(w, resp.Message, resp.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	prettyJSON, err := json.MarshalIndent(resp.Content, "", "  ")
	if err != nil {
		http.Error(w, "Failed to generate JSON output", http.StatusInternalServerError)
		return
	}

	w.Write(prettyJSON)
}
