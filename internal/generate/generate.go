package generate

import (
	"database/sql"
	"encoding/json"
	"go-test/internal/data"
	"go-test/internal/llm"
	"go-test/internal/prompts"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) ([]map[string]interface{}, error) { // should this return a status code?
	model, close, err := llm.NewClient(r.Context(), llm.Gemini_1_5)
	if err != nil {
		return nil, err
	}
	defer close()

	inputQuery := r.URL.Query().Get("query")
	if inputQuery == "" {
		return []map[string]interface{}{{"no query": true}}, nil // TODO: figure out better error handling
	}

	db, err := data.CreateDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	processOutput := func(output string) (string, error) {
		parsed, err := parseResponse(output)
		if err != nil {
			return "", err
		} else if parsed.Err {
			return "", err // TODO: fix this
		}
		return parsed.Sql, nil
	}

	getPrompts := []llm.LLMStep{
		{
			GetPrompt: func(lastOutput string) (string, error) {
				return prompts.GetChatPrompt(&prompts.ChatContext{Query: lastOutput})
			},
			ProcessOutput: processOutput,
		}, {
			GetPrompt: func(lastOutput string) (string, error) {
				return prompts.GetFixQueryPrompt(&prompts.ChatContext{Query: lastOutput})
			},
			ProcessOutput: processOutput,
		},
	}

	resp, err := model.GenerateSequence(r.Context(), getPrompts, inputQuery)
	if err != nil {
		return nil, err
	}

	queryResult, err := selectQueryDB(db, resp)
	if err != nil {
		return nil, err
	}

	return queryResult, nil
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

func selectQueryDB(db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Store a map of column id to value for each row
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		row := make([]interface{}, len(columns))
		for i := range columns {
			row[i] = new(interface{})
		}
		if err := rows.Scan(row...); err != nil {
			return nil, err
		}

		columnToValue := make(map[string]interface{})
		for i, colName := range columns {
			val := row[i].(*interface{})
			columnToValue[colName] = *val
		}
		result = append(result, columnToValue)
	}

	return result, nil
}
