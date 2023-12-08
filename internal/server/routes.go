package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ReviewPayload struct {
	ReviewText string `json:"reviewText"`
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)
	r.Get("/api/movies", s.GetAllMoviesHandler)
	r.Get("/api/movies/{title}", s.GetMovieHandler)
	r.Get("/api/movies/reviews/{title}", s.GetReviewsHandler)
	r.Post("/api/movies/add-review/{title}/{stars}", s.AddReviewHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) GetAllMoviesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.db.GetMovies())
}

func (s *Server) GetMovieHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	getMovie, err := s.db.GetMovie(title)
	if err != nil {
		log.Fatalf("Failed to get movie. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getMovie)
}

func (s *Server) GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	getReviews, err := s.db.ShowReview(title)
	if err != nil {
		if errors.Is(err, errors.New("no reviews")) {
			http.Error(w, "No reviews found", http.StatusNotFound)
			return
		} else {
			log.Printf("Failed to get review. Err: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getReviews)
}

func (s *Server) AddReviewHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	stars := chi.URLParam(r, "stars")
	var payload ReviewPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reviewText := payload.ReviewText

	starsNum, err := strconv.Atoi(stars)
	if err != nil {
		log.Fatalf("failed converting ratingstars string to int")
	}
	s.db.AddReview(title, starsNum, reviewText)
}
