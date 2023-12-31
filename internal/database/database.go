package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
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
	GetDirectors() []Director
	GetDirector(id string) (Director, error)
	GetActors() []Actor
	GetActor(id string) (Actor, error)
	ShowReview(url string) ([]Review, error)
	AddReview(url string, stars int, reviewText string, userName string)
	AuthenticateUser(username string, password string) (User, string, string)
	RegisterUser(username string, password string, email string) (string, error)
	//GetUserData(id int) (User, error)
	ToggleWatchlist(movieID, userID int) error
	ToggleLiked(movieID, userID int) error
	GetMoviesByDirectorID(directorID int) ([]DirectedMovie, error)
	GetMoviesByActorID(actorID int) ([]ActedMovie, error)
	GetStaffByMovieID(movieID int) ([]StaffMember, error)
	GetUserID(username string) (int, error)
	GetWatchlistStatus(movieID int, username string) string
	GetLikedStatus(movieID int, username string) string
}

type StaffMember struct {
	ID     int    `db:"staff_id" json:"staff_id"`
	MTitle string `json:"movie_title"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	TypeID int    `json:"type_id"`
}

type ActedMovie struct {
	Movie
	ActorID     int    `db:"actor_id" json:"actor_id"`
	ActorName   string `json:"ActorName"`
	ActorDob    string `json:"DateOfBirth"`
	Nationality string `json:"Nationality"`
}

type DirectedMovie struct {
	Movie
	DirectorID   int    `db:"director_id" json:"director_id"`
	DirectorName string `json:"DirectorName"`
	DirectorDob  string `json:"DateOfBirth"`
	Nationality  string `json:"Nationality"`
}

type Actor struct {
	ID          int    `db:"actor_id" json:"actor_id"`
	Name        string `json:"ActorName"`
	Dob         string `json:"DateOfBirth"`
	Nationality string `json:"Nationality"`
}

type Director struct {
	ID          int    `db:"director_id" json:"director_id"`
	Name        string `json:"DirectorName"`
	Dob         string `json:"DateOfBirth"`
	Nationality string `json:"Nationality"`
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

// var (
// 	dbname   = os.Getenv("DB_DATABASE")
// 	password = os.Getenv("DB_PASSWORD")
// 	username = os.Getenv("DB_USERNAME")
// 	port     = os.Getenv("DB_PORT")
// 	host     = os.Getenv("DB_HOST")
// )

func New() Service {
	// Opening a driver typically will not attempt to connect to the database.
	// db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
    db, err := sql.Open("mysql", "lab2324omada7:lab2324omada7@tcp(godockerDB)/lab2324omada7_tainia")
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

func (s *service) GetUserID(username string) (int, error) {
	selectDataQuery := fmt.Sprintf("SELECT user_id FROM USER where Username=%q", username)

	var userID int
	err := s.db.QueryRow(selectDataQuery).Scan(&userID)
	if err != nil {
		return -1, errors.New("User id doesn't exist") 
	}

	return userID, nil
}

func (s *service) GetMovies() []Movie {
	selectDataQuery := "SELECT * FROM MOVIE"

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

func (s *service) GetLikedStatus(movieID int, username string) string {
	selectDataQuery := "SELECT EXISTS (SELECT 1 FROM LIKES WHERE movie_id = ? AND user_id = ?) AS likes_status"

	var likesStatus string
    foo, err := s.GetUserID(username)
    if err != nil {
        log.Print(err)
    }

	err = s.db.QueryRow(selectDataQuery, movieID, foo).Scan(&likesStatus)
	if err != nil {
		log.Printf("Error checking likes status: %v", err)
		return "error"
	}
	fmt.Println(likesStatus)

	return likesStatus
}

func (s *service) GetWatchlistStatus(movieID int, username string) string {
	selectDataQuery := "SELECT EXISTS (SELECT 1 FROM ADDS_TO_WATCHLIST WHERE movie_id = ? AND user_id = ?) AS watchlist_status"

	var watchlistStatus string
    foo, err := s.GetUserID(username)
    if err != nil {
        log.Print(err)
    }

	err = s.db.QueryRow(selectDataQuery, movieID, foo).Scan(&watchlistStatus)
	if err != nil {
		log.Printf("Error checking watchlist status: %v", err)
		return "error"
	}
	return watchlistStatus
}

func (s *service) GetMovie(url string) (Movie, error) {
	modifiedTitle := strings.ReplaceAll(url, "-", " ")
	selectDataQuery := fmt.Sprintf("SELECT * FROM MOVIE WHERE Title=%q", modifiedTitle)

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

func (s *service) GetActors() []Actor {
	selectDataQuery := "SELECT * FROM ACTOR"

	rows, err := s.db.Query(selectDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var actors []Actor

	for rows.Next() {
		var actor Actor
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Dob, &actor.Nationality)
		if err != nil {
			panic(err.Error())
		}
		actors = append(actors, actor)
	}

	fmt.Println(actors)

	return actors
}

func (s *service) GetDirectors() []Director {
	selectDataQuery := "SELECT * FROM DIRECTOR"

	rows, err := s.db.Query(selectDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var directors []Director

	for rows.Next() {
		var director Director
		err := rows.Scan(&director.ID, &director.Name, &director.Dob, &director.Nationality)
		if err != nil {
			panic(err.Error())
		}
		directors = append(directors, director)
	}

	fmt.Println(directors)

	return directors
}

func (s *service) GetDirector(id string) (Director, error) {
	idNum, _ := strconv.Atoi(id)
	selectDataQuery := fmt.Sprintf("SELECT * FROM DIRECTOR WHERE director_id=%d", idNum)

	directorRow, err := s.db.Query(selectDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer directorRow.Close()

	if directorRow.Next() {
		var director Director
		err := directorRow.Scan(&director.ID, &director.Name, &director.Dob, &director.Nationality)
		if err != nil {
			return Director{}, err
		}
		return director, nil
	}

	fmt.Println(Director{})

	return Director{}, errors.New("director not found")
}

func (s *service) GetActor(id string) (Actor, error) {
	idNum, _ := strconv.Atoi(id)
	selectDataQuery := fmt.Sprintf("SELECT * FROM ACTOR WHERE actor_id=%d", idNum)

	actorRow, err := s.db.Query(selectDataQuery)
	if err != nil {
		panic(err.Error())
	}
	defer actorRow.Close()

	if actorRow.Next() {
		var actor Actor
		err := actorRow.Scan(&actor.ID, &actor.Name, &actor.Dob, &actor.Nationality)
		if err != nil {
			return Actor{}, err
		}
		return actor, nil
	}

	fmt.Println(Actor{})

	return Actor{}, errors.New("actor not found")
}

func (s *service) GetStaffByMovieID(movieID int) ([]StaffMember, error) {
	query := `
     SELECT
         M.movie_id,
         M.Title AS movie_title,
         A.actor_id,
         A.ActorName AS actor_name,
         'Actor' AS role
     FROM
         MOVIE M
         JOIN ACTED ACT ON M.movie_id = ACT.movie_id
         JOIN ACTOR A ON ACT.actor_id = A.actor_id
     WHERE
         M.movie_id = %d
     UNION
     SELECT
         M.movie_id,
         M.Title AS movie_title,
         D.director_id,
         D.DirectorName AS director_name,
         'Director' AS role
     FROM
         MOVIE M
         JOIN DIRECTED DIR ON M.movie_id = DIR.movie_id
         JOIN DIRECTOR D ON DIR.director_id = D.director_id
     WHERE
         M.movie_id = %d;
     `

	foo := fmt.Sprintf(query, movieID, movieID)

	rows, err := s.db.Query(foo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var staff []StaffMember
	for rows.Next() {
		var member StaffMember
		err := rows.Scan(
			&member.ID,
			&member.MTitle,
			&member.TypeID,
			&member.Name,
			&member.Role,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		staff = append(staff, member)
	}

	return staff, nil
}

func (s *service) GetMoviesByDirectorID(directorID int) ([]DirectedMovie, error) {
	query := `
		SELECT M.movie_id, M.Title, M.ReleaseDate, M.Genre, M.AvgRating, D.director_id, DIR.DateOfBirth, DIR.DirectorName, DIR.Nationality
		FROM MOVIE M
		JOIN DIRECTED D ON M.movie_id = D.movie_id
		JOIN DIRECTOR DIR ON D.director_id = DIR.director_id
		WHERE D.director_id = ?`

	rows, err := s.db.Query(query, directorID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var movies []DirectedMovie
	for rows.Next() {
		var movie DirectedMovie
		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.Genre,
			&movie.AvgRating,
			&movie.DirectorID,
			&movie.DirectorDob,
			&movie.DirectorName,
			&movie.Nationality,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *service) GetMoviesByActorID(actorID int) ([]ActedMovie, error) {
	query := `
		SELECT M.movie_id, M.Title, M.ReleaseDate, M.Genre, M.AvgRating, ACT.actor_id, ACT.DateOfBirth, ACT.ActorName, ACT.Nationality
		FROM MOVIE M
		JOIN ACTED A ON M.movie_id = A.movie_id
		JOIN ACTOR ACT ON A.actor_id = ACT.actor_id
		WHERE A.actor_id = ?`

	rows, err := s.db.Query(query, actorID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var movies []ActedMovie
	for rows.Next() {
		var movie ActedMovie
		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.Genre,
			&movie.AvgRating,
			&movie.ActorID,
			&movie.ActorDob,
			&movie.ActorName,
			&movie.Nationality,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// func (s *service) GetUserData(id int) (User, error) {
// 	selectDataQuery := fmt.Sprintf("SELECT * FROM USER WHERE user_id=%d", id)
// 
// 	userRow, err := s.db.Query(selectDataQuery)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer userRow.Close()
// 
// 	if userRow.Next() {
// 		var user User
// 		err := userRow.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
// 		if err != nil {
// 			return User{}, err
// 		}
// 		return user, nil
// 	}
// 
// 	fmt.Println(User{})
// 
// 	return User{}, errors.New("user not found")
// }

func (s *service) ShowReview(url string) ([]Review, error) {
	modifiedTitle := strings.ReplaceAll(url, "-", " ")
	selectDataQuery := fmt.Sprintf("SELECT * FROM MOVIE WHERE Title=%q", modifiedTitle)

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

	reviewDataQuery := fmt.Sprintf("SELECT * FROM REVIEW WHERE movie_id=%d", movie.Id)

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

func (s *service) AddReview(url string, stars int, reviewText string, username string) {
	modifiedTitle := strings.ReplaceAll(url, "-", " ")

	selectDataQuery := fmt.Sprintf("SELECT * FROM MOVIE WHERE Title=%q", modifiedTitle)
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

	getUserIDQuery := fmt.Sprintf("SELECT user_id FROM USER WHERE Username=%q", username)
	var userID int
	err = s.db.QueryRow(getUserIDQuery).Scan(&userID)
	if err != nil {
		panic(err.Error())
	}

	existingReviewQuery := "SELECT R.review_id FROM WROTE W JOIN REVIEW R ON W.review_id = R.review_id WHERE W.user_id = ? AND R.movie_id = ?"
	rows, err := s.db.Query(existingReviewQuery, userID, movie.Id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var existingReviewIDs []int
	for rows.Next() {
		var reviewID int
		if err := rows.Scan(&reviewID); err != nil {
			panic(err.Error())
		}
		existingReviewIDs = append(existingReviewIDs, reviewID)
	}

	currentTime := time.Now()
	dateToday := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	if len(existingReviewIDs) > 0 {
		updateReviewQuery := "UPDATE REVIEW SET ReviewText = ?, RatingStars = ?, DatePosted = ? WHERE review_id = ?"
		_, err = s.db.Exec(updateReviewQuery, reviewText, stars, dateToday, existingReviewIDs[len(existingReviewIDs)-1])
		if err != nil {
			panic(err.Error())
		}
	} else {
		insertReviewQuery := fmt.Sprintf("INSERT INTO REVIEW (ReviewText, RatingStars, DatePosted, movie_id) VALUES (%q, %d, '%s', %d);", reviewText, stars, dateToday, movie.Id)

		_, err = s.db.Exec(insertReviewQuery)
		if err != nil {
			panic(err.Error())
		}

		getLastReviewIDQuery := "SELECT LAST_INSERT_ID()"
		var lastReviewID int
		err = s.db.QueryRow(getLastReviewIDQuery).Scan(&lastReviewID)
		if err != nil {
			panic(err.Error())
		}

		insertWroteQuery := "INSERT INTO WROTE (review_id, user_id) VALUES (?, ?)"
		_, err = s.db.Exec(insertWroteQuery, lastReviewID, userID)
		if err != nil {
			panic(err.Error())
		}
	}

	updateAvgRatingQuery := fmt.Sprintf("UPDATE MOVIE SET AvgRating = (SELECT AVG(RatingStars) FROM REVIEW WHERE movie_id = %d) WHERE movie_id = %d", movie.Id, movie.Id)
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

	insertUserQuery := fmt.Sprintf("INSERT INTO USER (Username, Password, Email) VALUES (%q, %q, %q)", username, hashedPassword, email)
	_, err = s.db.Exec(insertUserQuery)
	if err != nil {
		return "", err
	}

	getUserIdQuery := fmt.Sprintf("SELECT user_id FROM USER WHERE Username=%q", username)
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

func (s *service) AuthenticateUser(username string, password string) (User, string, string) {
	selectUserQuery := fmt.Sprintf("SELECT * FROM USER WHERE Username=%q", username)
	userRow, err := s.db.Query(selectUserQuery)
	if err != nil {
		return User{}, "", "database error"
	}
	defer userRow.Close()

	var user User
	if userRow.Next() {
		err := userRow.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return User{}, "", "database error"
		}
	} else {
		return User{}, "", "usernotfound"
	}

	if comparePasswords(user.Password, password) {
		log.Printf("Authentication successful for user ID: %d, username: %s", user.ID, user.Username)
		token, err := createToken(user.ID)
		if err != nil {
			log.Println("Error creating token:", err)
			return User{}, "", "token creation error"
		}

		return user, token, ""
	} else {
		return User{}, "", "passNoMatch"
	}
}

func (s *service) ToggleWatchlist(movieID, userID int) error {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM ADDS_TO_WATCHLIST WHERE movie_id = ? AND user_id = ?)", movieID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = s.db.Exec("DELETE FROM ADDS_TO_WATCHLIST WHERE movie_id = ? AND user_id = ?", movieID, userID)
		if err != nil {
			return err
		}
	} else {
		_, err = s.db.Exec("INSERT INTO ADDS_TO_WATCHLIST (movie_id, user_id, DateAdded) VALUES (?, ?, ?)", movieID, userID, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) ToggleLiked(movieID, userID int) error {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM LIKES WHERE movie_id = ? AND user_id = ?)", movieID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = s.db.Exec("DELETE FROM LIKES WHERE movie_id = ? AND user_id = ?", movieID, userID)
		if err != nil {
			return err
		}
	} else {
		_, err = s.db.Exec("INSERT INTO LIKES (movie_id, user_id, DateAdded) VALUES (?, ?, ?)", movieID, userID, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
