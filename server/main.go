package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dtsmith94/shared-expenses-tracker/server/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on the port 8080...")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
