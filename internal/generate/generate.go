package generate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-test/internal/data"
	"go-test/internal/llm"
	"go-test/internal/prompts"
	"net/http"
)

type GenerateSqlResponse struct {
	Status  int
	Content []map[string]interface{}
	Message string
}

var ErrCreatingSqlQuery = errors.New("error creating sql query")

func GenerateHandler(r *http.Request) *GenerateSqlResponse {
	model, close, err := llm.NewClient(r.Context(), llm.Gemini_1_5)
	if err != nil {
		return &GenerateSqlResponse{
			Content: nil,
			Message: "Error creating model",
			Status:  http.StatusInternalServerError,
		}
	}
	defer close()

	inputQuery := r.URL.Query().Get("query")
	if inputQuery == "" {
		return &GenerateSqlResponse{
			Status:  http.StatusBadRequest,
			Content: nil,
			Message: "No query provided",
		}
	}

	db, err := data.CreateDB()
	if err != nil {
		return &GenerateSqlResponse{
			Content: nil,
			Message: "Error creating database",
			Status:  http.StatusInternalServerError,
		}
	}
	defer db.Close()

	processOutput := func(output string) (string, error) {
		parsed, err := parseResponse(output)
		if err != nil {
			return "", err
		} else if parsed.Err {
			return "", ErrCreatingSqlQuery
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
		// TODO: log
		return &GenerateSqlResponse{
			Content: nil,
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error while generating SQL",
		}
	}

	queryResult, err := selectQueryDB(db, resp)
	if err != nil {
		return &GenerateSqlResponse{
			Content: nil,
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error while running generated SQL",
		}
	}

	return &GenerateSqlResponse{
		Status:  http.StatusOK,
		Message: "",
		Content: queryResult,
	}
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
