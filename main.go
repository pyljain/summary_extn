package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"summariser/llm"
)

func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests (OPTIONS method)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

type Request struct {
	Content string `json:"content"`
}

type Response struct {
	Summary string `json:"summary"`
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/summarize", func(w http.ResponseWriter, r *http.Request) {
		aiReq := Request{}
		log.Printf("Handling request to AI")
		err := json.NewDecoder(r.Body).Decode(&aiReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		prompt := fmt.Sprintf("Please summarise the following content as concisely as possible. Organise your summary into paragraphs and provide sections that are relevant to the content\n```%s```", aiReq.Content)

		resp, err := llm.CallOpenAI(prompt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rp := Response{
			Summary: resp,
		}

		err = json.NewEncoder(w).Encode(&rp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	err := http.ListenAndServe(":9001", allowCORS(mux))
	if err != nil {
		log.Fatalf("Error starting server %s", err)
	}
}
