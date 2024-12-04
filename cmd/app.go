package main

import (
	"context"
	"fmt"
	"log"

	"go-test/internal/llm"
	"go-test/internal/prompts"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	// client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Close()

	// model := client.GenerativeModel("gemini-1.5-pro-latest")

	// resp, err := model.GenerateContent(ctx, genai.Text("List a few popular cookie recipes using this JSON schema."))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// printResponse(resp)

	model, err := llm.NewClient(ctx, llm.Gemini_1_5)
	if err != nil {
		log.Fatal(err)
	}

	prompt, err := prompts.GetTestPrompt(&prompts.TestPromptContext{
		Food: "chicken nuggets",
	})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := model.Generate(ctx, prompt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)

	err = model.Close()
	if err != nil {
		log.Fatal(err)
	}

	// router := httprouter.New()
	// router.GET("/", routes.Index)
	// router.GET("/hello/:name", routes.Hello)

	// log.Fatal(http.ListenAndServe(":8080", router))
}
