package main

import (
	"fmt"
	"log"
	"net/http"

	"go-test/internal/routes"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := httprouter.New()
	router.GET("/hello", routes.Hello)
	router.GET("/generate", routes.GenerateSQL)

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
