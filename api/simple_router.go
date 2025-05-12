package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	var error map[string]string
	switch status {
	case http.StatusMethodNotAllowed:
		error = map[string]string{
			"error": fmt.Sprintf("%s Method Not Allowed", r.Method),
		}

	case http.StatusNotFound:
		error = map[string]string{
			"error": fmt.Sprintf("%s Path does not exist", r.URL.Path),
		}
	case http.StatusBadRequest:
		error = map[string]string{
			"error": "Invalid JSON",
		}
	default:
		error = map[string]string{
			"error": "Unexpected Error",
		}
	}
	json.NewEncoder(w).Encode(error)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "pong!")
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	type Input struct {
		Name string `json:"name"`
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	response := map[string]string{
		"message": fmt.Sprintf("Welcome, %s", input.Name),
	}
	json.NewEncoder(w).Encode(response)

}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"message": "This is an about Page.",
	}
	json.NewEncoder(w).Encode(data)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"message": "This is a contact Page.",
	}
	json.NewEncoder(w).Encode(data)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"message": "This is a product Page.",
	}
	json.NewEncoder(w).Encode(data)
}

type customMux struct{}

func (m *customMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/about":
		aboutHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/ping":
		pingHandler(w, r)
	case "/product":
		productHandler(w, r)
	case "/welcome":
		welcomeHandler(w, r)
	default:
		errorHandler(w, r, http.StatusNotFound)
	}
}

func main() {
	PORT_NUMER := 8080
	log.Printf("Starting Server at http://localhost:%d", PORT_NUMER)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT_NUMER), &customMux{})
	if err != nil {
		log.Fatalf("Error Starting Server: %v", err)
	}
}
