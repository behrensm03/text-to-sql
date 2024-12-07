package prompts

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

type ChatContext struct {
	Query string
}

//go:embed templates/*
var templates embed.FS

func GetChatPrompt(ctx *ChatContext) (string, error) {
	return executeTemplate(ctx, "chat-min.tmpl")
}

func GetFixQueryPrompt(ctx *ChatContext) (string, error) {
	return executeTemplate(ctx, "fix.tmpl")
}

func executeTemplate[T any](data T, file string) (string, error) {
	tmpl, err := template.ParseFS(templates, "templates/"+file)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return result.String(), nil
}
