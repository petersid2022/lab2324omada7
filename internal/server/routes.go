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
)

var JWT_SECRET = []byte(fmt.Sprint(os.Getenv("KEY")))

type LikedPayload struct {
	MovieID  string `json:"movieId"`
	Username string `json:"userName"`
}

type WatchlistPayload struct {
	MovieID  string `json:"movieId"`
	Username string `json:"userName"`
}

type UserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReviewPayload struct {
	ReviewText   string `json:"reviewText"`
	UserNametext string `json:"userName"`
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", s.healthHandler)
	r.Get("/api/directors", s.GetAllDirectorsHandler)
	r.Get("/api/movies/staff/{name}", s.GetAllMovieStaffHandler)
	r.Get("/api/directors/{name}", s.GetDirectorHandler)
	r.Get("/api/actors", s.GetAllActorsHandler)
	r.Get("/api/actors/{name}", s.GetActorHandler)
	r.Get("/api/movies", s.GetAllMoviesHandler)
	r.Get("/api/movies/{title}", s.GetMovieHandler)
	r.Get("/api/movies/reviews/{title}", s.GetReviewsHandler)
	r.Get("/userdata/{id}", s.UserDataHandler)
	r.Get("/api/directors/{id}", s.DirectedHandler)
	r.Get("/api/actors/{id}", s.ActedHandler)
	r.Post("/api/movies/add-review/{title}/{stars}", s.AddReviewHandler)
	r.Post("/create-account", s.CreateAccountHandler)
	r.Post("/login", s.LoginHandler)
	r.Post("/api/watchlist", s.ToggleWatchlistHandler)
	r.Get("/watchlistStatus/{movie_id}/{username}", s.GetWatchlistHandler)
	r.Post("/api/liked", s.ToggleLikedHandler)
	r.Get("/likedStatus/{movie_id}/{username}", s.GetLikedHandler)

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

func (s *Server) GetAllDirectorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.db.GetDirectors())
}

func (s *Server) GetAllActorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.db.GetActors())
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

func (s *Server) GetDirectorHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	getDirector, err := s.db.GetDirector(name)
	if err != nil {
		log.Fatalf("Failed to get director. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getDirector)
}

func (s *Server) GetActorHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	getActor, err := s.db.GetActor(name)
	if err != nil {
		log.Fatalf("Failed to get actor. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getActor)
}

func (s *Server) GetAllMovieStaffHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	movie, err := s.db.GetMovie(name)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Failed to get movie. Err: %v", err)
	}
	getStaff, err := s.db.GetStaffByMovieID(movie.Id)
	if err != nil {
		log.Fatalf("Failed to get staff. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getStaff)
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
	userNameText := payload.UserNametext

	starsNum, err := strconv.Atoi(stars)
	if err != nil {
		log.Fatalf("failed converting ratingstars string to int")
	}
	s.db.AddReview(title, starsNum, reviewText, userNameText)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ok")
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
	user, token, errMsg := s.db.AuthenticateUser(username, password)
	if errMsg != "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": errMsg,
		})
		return
	}

	fmt.Println("User logged in successfully:", username)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data": map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
}

func (s *Server) UserDataHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idNum, _ := strconv.Atoi(id)
	fmt.Println(idNum)
	getUserdata, err := s.db.GetUserData(idNum)
	if err != nil {
		if errors.Is(err, errors.New("no user")) {
			http.Error(w, "No user found", http.StatusNotFound)
			return
		} else {
			log.Printf("Failed to get review. Err: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getUserdata)
}

func (s *Server) DirectedHandler(w http.ResponseWriter, r *http.Request) {
	directorIDStr := chi.URLParam(r, "id")
	directorID, _ := strconv.Atoi(directorIDStr)
	movies, err := s.db.GetMoviesByDirectorID(directorID)
	if err != nil {
		http.Error(w, "Error fetching directors", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(movies)
	json.NewEncoder(w).Encode(movies)
}

func (s *Server) ActedHandler(w http.ResponseWriter, r *http.Request) {
	actedIdStr := chi.URLParam(r, "id")
	actedID, _ := strconv.Atoi(actedIdStr)
	movies, err := s.db.GetMoviesByActorID(actedID)
	if err != nil {
		http.Error(w, "Error fetching actors", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(movies)
	json.NewEncoder(w).Encode(movies)
}

func (s *Server) ToggleWatchlistHandler(w http.ResponseWriter, r *http.Request) {
	var payload WatchlistPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := payload.Username
	userid := s.db.GetUserID(username)
	movieId := payload.MovieID
	movieIdNum, _ := s.db.GetMovie(movieId)
	err := s.db.ToggleWatchlist(movieIdNum.Id, userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("toggled watchlist successfully for user:", username, " and movie_id:", movieId)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (s *Server) ToggleLikedHandler(w http.ResponseWriter, r *http.Request) {
	var payload LikedPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := payload.Username
	userid := s.db.GetUserID(username)
	movieId := payload.MovieID
	movieIdNum, _ := s.db.GetMovie(movieId)
	err := s.db.ToggleLiked(movieIdNum.Id, userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("toggled liked successfully for user:", username, " and movie_id:", movieId)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (s *Server) GetWatchlistHandler(w http.ResponseWriter, r *http.Request) {
	movieIdStr := chi.URLParam(r, "movie_id")
	movieID, _ := strconv.Atoi(movieIdStr)
	username := chi.URLParam(r, "username")
	status := s.db.GetWatchlistStatus(movieID, username)
	//fmt.Println("watchlist status for user:", username, " and movie_id:", movieID, " == ", status)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"data":   status,
	})
}

func (s *Server) GetLikedHandler(w http.ResponseWriter, r *http.Request) {
	movieIdStr := chi.URLParam(r, "movie_id")
	movieID, _ := strconv.Atoi(movieIdStr)
	username := chi.URLParam(r, "username")
	status := s.db.GetLikedStatus(movieID, username)
	//fmt.Println("watchlist status for user:", username, " and movie_id:", movieID, " == ", status)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"data":   status,
	})
}

//func (s *Server) UserDataHandler(w http.ResponseWriter, r *http.Request) {
//	decoder := json.NewDecoder(r.Body)
//	var requestData struct {
//		Token string `json:"token"`
//	}
//	err := decoder.Decode(&requestData)
//	if err != nil {
//		http.Error(w, "Invalid request body", http.StatusBadRequest)
//		return
//	}
//
//	fmt.Println("Received Token:", requestData.Token)
//
//	token, err := jwt.Parse(requestData.Token, func(token *jwt.Token) (interface{}, error) {
//		return JWT_SECRET, nil
//	})
//
//	if err != nil {
//		fmt.Println("Error parsing token:", err)
//		http.Error(w, "Token verification failed", http.StatusUnauthorized)
//		return
//	}
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		var userEmail, userUsername string
//
//		if emailClaim, ok := claims["email"].(string); ok {
//			userEmail = emailClaim
//		}
//
//		if usernameClaim, ok := claims["username"].(string); ok {
//			userUsername = usernameClaim
//		}
//
//		userData := UserPayload{
//			Username: userUsername,
//			Email:    userEmail,
//			Password: "",
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(map[string]interface{}{
//			"status": "ok",
//			"data":   userData,
//		})
//	} else {
//		http.Error(w, "Invalid token", http.StatusUnauthorized)
//	}
//
//}
