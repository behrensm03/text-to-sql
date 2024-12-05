package generate

import (
	"encoding/json"
	"fmt"
	"go-test/internal/llm"
	"go-test/internal/prompts"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (string, error) { // should this return a status code?
	model, err := llm.NewClient(r.Context(), llm.Gemini_1_5)
	if err != nil {
		return "", err
	}

	// TODO: get query from url
	prompt, err := prompts.GetChatPrompt(&prompts.ChatContext{
		Query: "Show me all the customer names",
	})
	if err != nil {
		return "", err
	}

	resp, err := model.Generate(r.Context(), prompt)
	if err != nil {
		return "", err
	}

	result, err := parseResponse(resp)
	if err != nil {
		return "", err
	}

	err = model.Close()
	if err != nil {
		return "", err
	}

	fmt.Println(result)

	return "", nil
}

type promptOutput struct {
	Sql string `json:"sql"`
	Err bool   `json:"error"`
}

func parseResponse(resp string) (*promptOutput, error) {
	var result promptOutput
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
