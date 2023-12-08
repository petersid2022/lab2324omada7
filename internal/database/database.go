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
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	GetMovies() []Movie
	GetMovie(url string) (Movie, error)
	ShowReview(url string) ([]Review, error)
	AddReview(url string, stars int, reviewText string)
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

	reviewDataQuery := fmt.Sprintf("SELECT * FROM review WHERE movie_id=%d", movie.Id)

	reviewRow, err := s.db.Query(reviewDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer reviewRow.Close()

	currentTime := time.Now()
	dateToday := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	insertDataQuery := fmt.Sprintf("INSERT INTO review (ReviewText, RatingStars, DatePosted, movie_id) VALUES (%q, %d, '%s', %d);", reviewText, stars, dateToday, movie.Id)

	_, err = s.db.Exec(insertDataQuery)
	if err != nil {
		panic(err.Error())
	}
}
