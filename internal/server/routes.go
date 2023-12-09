package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte(fmt.Sprint(os.Getenv("KEY")))

type UserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReviewPayload struct {
	ReviewText string `json:"reviewText"`
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", s.healthHandler)
	r.Get("/api/movies", s.GetAllMoviesHandler)
	r.Get("/api/movies/{title}", s.GetMovieHandler)
	r.Get("/api/movies/reviews/{title}", s.GetReviewsHandler)
	r.Post("/api/movies/add-review/{title}/{stars}", s.AddReviewHandler)
	r.Post("/create-account", s.CreateAccountHandler)
	r.Post("/login", s.LoginHandler)
	r.Post("/userdata", s.UserDataHandler)

	return r
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

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var payload UserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := payload.Username
	email := payload.Email
	password := payload.Password
	register, err := s.db.RegisterUser(username, password, email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("User registered successfully:", username, email)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data":   register,
	})
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload UserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := payload.Username
	password := payload.Password
	out, err := s.db.AuthenticateUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("User logged in successfully:", username)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data":   out,
	})
}

func (s *Server) UserDataHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestData struct {
		Token string `json:"token"`
	}
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Received Token:", requestData.Token)

	token, err := jwt.Parse(requestData.Token, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		http.Error(w, "Token verification failed", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var userEmail, userUsername string

		if emailClaim, ok := claims["email"].(string); ok {
			userEmail = emailClaim
		}

		if usernameClaim, ok := claims["username"].(string); ok {
			userUsername = usernameClaim
		}

		userData := UserPayload{
			Username: userUsername,
			Email:    userEmail,
			Password: "",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
			"data":   userData,
		})
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}

}
