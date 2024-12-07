package llm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	ErrMissingContent      = errors.New("llm missing content")
	ErrContentTypeMismatch = errors.New("llm response content type mismatch")
)

type GetPromptFunction func(lastOutput string) (string, error)
type ProcessOutputFunction func(output string) (string, error)
type LLMStep struct {
	GetPrompt     GetPromptFunction
	ProcessOutput ProcessOutputFunction
}

type Model interface {
	Generate(ctx context.Context, prompt string) (string, error)
	GenerateSequence(ctx context.Context, steps []LLMStep, initialContext string) (string, error)
	Close() error
}

type ModelType string

const (
	Gemini_1_5 ModelType = "gemini-1-5"
)

type Gemini struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

type CloseFunc func()

func NewClient(ctx context.Context, modelType ModelType) (Model, CloseFunc, error) {
	switch modelType {
	case Gemini_1_5:
		gemini, err := newGemini(ctx)
		return gemini, func() {
			err := gemini.Close()
			if err != nil {
				log.Fatal(err)
			}
		}, err
	default:
		return nil, nil, fmt.Errorf("unrecognized model")
	}
}

func newGemini(ctx context.Context) (*Gemini, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	geminiModel := client.GenerativeModel("gemini-1.5-flash")
	// Ask the model to respond with JSON.
	geminiModel.ResponseMIMEType = "application/json"

	// Specify the schema.
	geminiModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"sql": {
				Type:        genai.TypeString,
				Description: "The sql string you generate",
			},
			"error": {
				Type:        genai.TypeBoolean,
				Description: "True if you are unable to generate a sql query, or false otherwise.",
			},
		},
	}

	model := &Gemini{
		client: client,
		model:  geminiModel,
	}

	return model, nil
}

func (g *Gemini) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	content, err := parseGeminiResponse(resp)
	if err != nil {
		return "", fmt.Errorf("failed to parse gemini response: %w", err)
	}

	return content, nil
}

func (g *Gemini) GenerateSequence(ctx context.Context, steps []LLMStep, initialContext string) (string, error) {
	currentOutput := initialContext
	for _, step := range steps {
		prompt, err := step.GetPrompt(currentOutput)
		if err != nil {
			return "", err
		}
		output, err := g.Generate(ctx, prompt)
		if err != nil {
			return "", err
		}
		parsed, err := step.ProcessOutput(output)
		if err != nil {
			return "", err
		}
		currentOutput = parsed
	}

	return currentOutput, nil
}

func (g *Gemini) Close() error {
	if err := g.client.Close(); err != nil {
		return fmt.Errorf("error closing gemini client: %w", err)
	}

	return nil
}

func parseGeminiResponse(resp *genai.GenerateContentResponse) (string, error) {
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", ErrMissingContent
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", ErrMissingContent
	}

	txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", ErrContentTypeMismatch
	}

	return string(txt), nil
}
