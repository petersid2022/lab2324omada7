package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte(fmt.Sprint(os.Getenv("KEY")))

type Service interface {
	Health() map[string]string
	GetMovies() []Movie
	GetMovie(url string) (Movie, error)
	ShowReview(url string) ([]Review, error)
	AddReview(url string, stars int, reviewText string)
	AuthenticateUser(username string, password string) (string, error)
	RegisterUser(username string, password string, email string) (string, error)
}

type User struct {
	ID       int    `json:"user_id"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
}

type Movie struct {
	Id          int     `json:"movie_id"`
	Title       string  `json:"Title"`
	ReleaseDate string  `json:"ReleaseDate"`
	Genre       string  `json:"Genre"`
	AvgRating   float64 `json:"AvgRating"`
}

type Review struct {
	Id         int    `json:"review_id"`
	Stars      int    `json:"RatingStars"`
	Review     string `json:"ReviewText"`
	DatePosted string `json:"DatePosted"`
	MovieId    string `json:"movie_id"`
}

type service struct {
	db *sql.DB
}

var (
	dbname   = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	// Opening a driver typically will not attempt to connect to the database.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	s := &service{db: db}
	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetMovies() []Movie {
	selectDataQuery := "SELECT * FROM movie"

	rows, err := s.db.Query(selectDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.Id, &movie.Title, &movie.ReleaseDate, &movie.Genre, &movie.AvgRating)
		if err != nil {
			panic(err.Error())
		}
		movies = append(movies, movie)
	}

	fmt.Println(movies)

	return movies
}

func (s *service) GetMovie(url string) (Movie, error) {
	modifiedTitle := strings.ReplaceAll(url, "-", " ")
	selectDataQuery := fmt.Sprintf("SELECT * FROM movie WHERE Title=%q", modifiedTitle)

	movieRow, err := s.db.Query(selectDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer movieRow.Close()

	if movieRow.Next() {
		var movie Movie
		err := movieRow.Scan(&movie.Id, &movie.Title, &movie.ReleaseDate, &movie.Genre, &movie.AvgRating)
		if err != nil {
			return Movie{}, err
		}
		return movie, nil
	}

	fmt.Println(Movie{})

	return Movie{}, errors.New("movie not found")
}

func (s *service) ShowReview(url string) ([]Review, error) {
	modifiedTitle := strings.ReplaceAll(url, "-", " ")
	selectDataQuery := fmt.Sprintf("SELECT * FROM movie WHERE Title=%q", modifiedTitle)

	movieRow, err := s.db.Query(selectDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer movieRow.Close()

	var movie Movie
	if movieRow.Next() {
		err := movieRow.Scan(&movie.Id, &movie.Title, &movie.ReleaseDate, &movie.Genre, &movie.AvgRating)
		if err != nil {
			panic(err.Error())
		}
	}

	reviewDataQuery := fmt.Sprintf("SELECT * FROM review WHERE movie_id=%d", movie.Id)

	reviewRow, err := s.db.Query(reviewDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer reviewRow.Close()

	var reviews []Review

	var ErrNoReviews = errors.New("no reviews")

	for reviewRow.Next() {
		var review Review
		err := reviewRow.Scan(&review.Id, &review.Review, &review.Stars, &review.DatePosted, &review.MovieId)
		if err != nil {
			log.Printf("Error scanning review row: %v", err)
		}
		reviews = append(reviews, review)
	}

	if len(reviews) == 0 {
		return nil, ErrNoReviews
	}

	return reviews, nil
}

func (s *service) AddReview(url string, stars int, reviewText string) {
	modifiedTitle := strings.ReplaceAll(url, "-", " ")
	selectDataQuery := fmt.Sprintf("SELECT * FROM movie WHERE Title=%q", modifiedTitle)

	movieRow, err := s.db.Query(selectDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer movieRow.Close()

	var movie Movie
	if movieRow.Next() {
		err := movieRow.Scan(&movie.Id, &movie.Title, &movie.ReleaseDate, &movie.Genre, &movie.AvgRating)
		if err != nil {
			panic(err.Error())
		}
	}

	currentTime := time.Now()
	dateToday := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	insertDataQuery := fmt.Sprintf("INSERT INTO review (ReviewText, RatingStars, DatePosted, movie_id) VALUES (%q, %d, '%s', %d);", reviewText, stars, dateToday, movie.Id)

	_, err = s.db.Exec(insertDataQuery)
	if err != nil {
		panic(err.Error())
	}

	reviewDataQuery := fmt.Sprintf("SELECT AVG(RatingStars) FROM review WHERE movie_id=%d", movie.Id)

	avgRatingRow, err := s.db.Query(reviewDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer avgRatingRow.Close()

	var avgRating float64
	if avgRatingRow.Next() {
		err := avgRatingRow.Scan(&avgRating)
		if err != nil {
			panic(err.Error())
		}
	}

	updateAvgRatingQuery := fmt.Sprintf("UPDATE movie SET AvgRating = %f WHERE movie_id = %d", avgRating, movie.Id)

	_, err = s.db.Exec(updateAvgRatingQuery)
	if err != nil {
		panic(err.Error())
	}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *service) RegisterUser(username string, password string, email string) (string, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", err
	}

	insertUserQuery := fmt.Sprintf("INSERT INTO user (Username, Password, Email) VALUES (%q, %q, %q)", username, hashedPassword, email)
	_, err = s.db.Exec(insertUserQuery)
	if err != nil {
		return "", err
	}

	getUserIdQuery := fmt.Sprintf("SELECT user_id FROM user WHERE Username=%q", username)
	var userID int
	err = s.db.QueryRow(getUserIdQuery).Scan(&userID)
	if err != nil {
		return "", err
	}

	token, err := createToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func comparePasswords(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}

func (s *service) AuthenticateUser(username string, password string) (string, error) {
	selectUserQuery := fmt.Sprintf("SELECT * FROM user WHERE Username=%q", username)
	userRow, err := s.db.Query(selectUserQuery)
	if err != nil {
		return "", err
	}
	defer userRow.Close()

	var user User
	if userRow.Next() {
		err := userRow.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("username not found")
	}

	if comparePasswords(user.Password, password) {
        log.Printf("Authentication successful for user ID: %d, username: %s", user.ID, user.Username)
		token, err := createToken(user.ID)
		if err != nil {
			log.Println("Error creating token:", err)
			return "", errors.New("Error creating token")
		}
		return token, nil
	} else {
		fmt.Printf("Hashed Password: %s, Provided Password: %s\n", user.Password, password)
	}

	return "", nil
}
