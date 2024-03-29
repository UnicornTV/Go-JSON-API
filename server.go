package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Define a custom port, which is nice both for dev and production
var port string = os.Getenv("PORT")

func main() {
	if port == "" {
		log.Println("no port name provided, using 3000")
		port = "3000"
	}

	// Since this is a very minimal API, we should have all requests come through the root handler
	http.HandleFunc("/", root)
	log.Printf("Listening for connections on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

/**
 * POST /
 *
 * Handle all requests to the root url.
 * This expects a POST request with a JSON request body
 * the payload should contain the values defined in the Input struct
 */
func root(w http.ResponseWriter, r *http.Request) {
	// Set the content type to JSON regardless of the data returned
	w.Header().Set("Content-Type", "application/json")

	// Try to decode the data provided
	var cereal Input
	err := json.NewDecoder(r.Body).Decode(&cereal)

	// If an error is encountered, let's say so with a nice handler
	if err != nil {
		handleError("Error parsing data", 500, w)
		return
	}

	// Build our recommendations with some simple math
	recommendations := Recommendations{
		Output{
			Description: fmt.Sprint(cereal.Name, " with skim milk"),
			Cereal:      cereal.Cups,
			Milk:        cereal.Cups * .75,
		},
		Output{
			Description: fmt.Sprint(cereal.Name, " with 2% milk"),
			Cereal:      cereal.Cups,
			Milk:        cereal.Cups * 1,
		},
		Output{
			Description: fmt.Sprint(cereal.Name, " with whole milk"),
			Cereal:      cereal.Cups,
			Milk:        cereal.Cups * 1.25,
		},
	}

	// Encode the recommendations array of structs into JSON
	json.NewEncoder(w).Encode(recommendations)
}

/**
 * Handle Errors
 *
 * Response codes are more than enough in most applications, however, I find that a nice JSON response body
 * with a more verbose description of the error and code are tremendously helpful to developers.
 */
func handleError(message string, code int, w http.ResponseWriter) {
	err := Error{Message: message, Status: code}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
