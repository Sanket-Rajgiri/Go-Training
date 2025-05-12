package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "pong!")
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		name = "Guest"
	}
	fmt.Fprintf(w, "Hello %s from Go!\n", name)
}

func byeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		name = "Guest"
	}
	fmt.Fprintf(w, "Bye %s!\n", name)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		name = "Guest"
	}
	data := map[string]string{
		"message": fmt.Sprintf("Hello %s", name),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
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
		"message": fmt.Sprintf("Hello %s", input.Name),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/bye", byeHandler)
	http.HandleFunc("/input", postHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/welcome", welcomeHandler)

	fmt.Println("Starting Server at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error Starting Server: %v", err)
	}
}
