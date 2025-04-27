package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricReset() {
	cfg.fileserverHits.Store(0)
}

func main() {
	addr := "localhost:8080"
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", apiCfg.fileserverHits.Load())))
	})
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.metricReset()
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset"))
	})
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		// set response content type
		w.Header().Set("Content-Type", "application/json")

		// Unmarshal request body into struct
		b := struct {
			Body string `json:"body"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, `{"error":"Something went wrong"}`, http.StatusBadRequest)
			return
		}
		if len(b.Body) > 140 {
			http.Error(w, `{"error":"Chirp too long"}`, http.StatusBadRequest)
			return
		}
		filteredBody := filterProfanity(b.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"cleaned_body": "%s"}`, filteredBody)))
	})
	log.Printf("Server listening on %s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func filterProfanity(text string) string {
	profanity := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	filteredText := []string{}
	for _, word := range strings.Fields(text) {
		w := strings.ToLower(word)
		if profanity[w] {
			filteredText = append(filteredText, "****")
		} else {
			filteredText = append(filteredText, word)
		}
	}
	return strings.Join(filteredText, " ")
}
