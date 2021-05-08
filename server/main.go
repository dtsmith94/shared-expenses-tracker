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
	log.Fatal(http.ListenAndServe(":8080", r))
}
