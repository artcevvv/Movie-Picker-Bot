package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// all queries in queries.go file

func InitDb() {
	var err error
	connectionString := "postgres://moviepickeruser:artcevvv!7437@localhost/moviepickerbot?sslmode=disable"
	db, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DATABASE CONNECTION IS SUCCESSFUL")

	// creates table for movies

	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Failed to create movies table: %v", err)
	}

	fmt.Printf("TABLE MOVIES CREATED\n\n")

	// creates users table

	_, err = db.Exec(queryForUsers)

	if err != nil {
		log.Printf("Failed to create users table: %v", err)
	}

	fmt.Printf("TABLE USERS CREATED\n\n")

	// creates series table

	_, err = db.Exec(queryForSeries)

	if err != nil {
		log.Printf("Failed to create series table: %v", err)
	}

	fmt.Printf("TABLE SERIES CREATED\n\n")
}

func addUser(telegramUserID int64, telegramUsername string) error {
	var count int

	err := db.QueryRow(isUserExists, telegramUserID).Scan(&count)

	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	if count == 0 {
		_, err := db.Exec(addUserQuery, telegramUserID, telegramUsername)

		if err != nil {
			return fmt.Errorf("error when adding user: %v", err)
		}

		fmt.Printf("User added successfully")
	}

	return nil
}

func addMovieHandler(telegramUsenameOwner string, telegramUserOwnerID int64, movieTitle, movieGenre string, telegramUserBoundedID *int64) error {
	if movieTitle == "" {
		return fmt.Errorf("movie title cannot be empty")
	}

	_, err := db.Exec(addMovieQuery, telegramUsenameOwner, telegramUserOwnerID, movieTitle, movieGenre, telegramUserBoundedID)

	if err != nil {
		return fmt.Errorf("failed to add movie: %s", err)
	}

	fmt.Println("Movie added successfully")
	return nil
}

func getMoviesHandler(telegramUsenameOwner string) ([]map[string]string, error) {
	rows, err := db.Query(getMoviesQuery, telegramUsenameOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %s", err)
	}

	defer rows.Close()

	var movies []map[string]string

	for rows.Next() {
		var movieTitle, movieGenre string
		err := rows.Scan(&movieTitle, &movieGenre)

		if err != nil {
			return nil, fmt.Errorf("failed to parse movies: %v", err)
		}

		movie := map[string]string{
			"title": movieTitle,
			"genre": movieGenre,
		}
		movies = append(movies, movie)
	}

	if len(movies) == 0 {
		return nil, fmt.Errorf("no movies found")
	}

	return movies, nil
}

func getMoviesByGenre(telegramUsenameOwner, movieGenre string) ([]map[string]string, error) {

	rows, err := db.Query(getByGenreQuery, telegramUsenameOwner, movieGenre)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %v", err)
	}

	defer rows.Close()

	var moviesByGenre []map[string]string
	for rows.Next() {
		var movieTitle string
		err := rows.Scan(&movieTitle)
		if err != nil {
			return nil, fmt.Errorf("failed to parse movies: %v", err)
		}

		movieByGenre := map[string]string{
			"title": movieTitle,
		}
		moviesByGenre = append(moviesByGenre, movieByGenre)
	}

	if len(moviesByGenre) == 0 {
		return nil, fmt.Errorf("no movies with this genre was found")
	}

	return moviesByGenre, err
}

func rmMovie(telegramUsenameOwner, movieTitle string) (string, error) {
	result, err := db.Exec(deleteQuery, telegramUsenameOwner, movieTitle)
	if err != nil {
		return "", fmt.Errorf("failed to remove movie: %v", err)
	}

	ra, err := result.RowsAffected()

	if err != nil {
		return "", fmt.Errorf("couldn't determine affected rows: %v", err)
	}

	if ra == 0 {
		return "No movies found for this user", nil
	}

	return fmt.Sprintf("Movie '%s' removed successfully.", movieTitle), nil
}
