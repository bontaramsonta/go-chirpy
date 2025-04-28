package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bontaramsonta/go-chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        int32     `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
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

func (cfg *apiConfig) createUser(ctx context.Context, email string) (User, error) {
	user, err := cfg.dbQueries.CreateUser(ctx, email)
	if err != nil {
		log.Println(err)
		return User{}, err
	}
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, nil
}

func (cfg *apiConfig) deleteAllUsers(ctx context.Context) error {
	return cfg.dbQueries.DeleteAllUsers(ctx)
}

func (cfg *apiConfig) createChirp(ctx context.Context, userID uuid.UUID, body string) (Chirp, error) {
	chirp, err := cfg.dbQueries.CreateChirp(ctx, database.CreateChirpParams{
		UserID: userID,
		Body:   body,
	})
	if err != nil {
		log.Println(err)
		return Chirp{}, err
	}
	return Chirp{
		ID:        chirp.ID,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}, nil
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	addr := "localhost:8080"
	apiCfg := &apiConfig{
		dbQueries: dbQueries,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", apiCfg.fileserverHits.Load())))
	})
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.metricReset()
		platform := os.Getenv("PLATFORM")
		if platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err := apiCfg.deleteAllUsers(r.Context())
		if err != nil {
			http.Error(w, `{"error":"delete all users failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset"))
	})
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		// set response content type
		w.Header().Set("Content-Type", "application/json")

		// Unmarshal request body into struct
		b := struct {
			Email string `json:"email"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, `{"error":"body json parse failed"}`, http.StatusBadRequest)
			return
		}
		user, err := apiCfg.createUser(r.Context(), b.Email)
		if err != nil {
			http.Error(w, `{"error":"create user failed"}`, http.StatusBadRequest)
			return
		}
		// unmarshall user to json and write response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		// set response content type
		w.Header().Set("Content-Type", "application/json")

		// Unmarshal request body into struct
		b := struct {
			UserID uuid.UUID `json:"user_id"`
			Body   string    `json:"body"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, `{"error":"body json parse failed"}`, http.StatusBadRequest)
			return
		}
		if len(b.Body) > 140 {
			http.Error(w, `{"error":"Chirp too long"}`, http.StatusBadRequest)
			return
		}
		filteredBody := filterProfanity(b.Body)
		chirp, err := apiCfg.createChirp(r.Context(), b.UserID, filteredBody)
		if err != nil {
			http.Error(w, `{"error":"create chirp failed"}`, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(chirp)
	})
	log.Printf("Server listening on %s\n", addr)
	err = http.ListenAndServe(addr, mux)
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
