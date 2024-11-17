package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDb() {
	var err error
	connectionString := "postgres://MoviePickerUser:7437@localhost/moviepickerbot?sslmode=disable"
	db, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DATABASE CONNECTION IS SUCCESSFUL")

	query := `
		CREATE TABLE IF NOT EXISTS movies (
			id SERIAL PRIMARY KEY,
			telegramUsenameOwner VARCHAR,
			telegramUserOwnerID BIGINT NOT NULL,
			movieTitle VARCHAR NOT NULL,
			movieGenre VARCHAR,
			telegramUserBoundedID BIGINT);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Failed to create movies table: %v", err)
	}

	fmt.Println("TABLE CREATED")
}

func addMovieHandler(telegramUsenameOwner string, telegramUserOwnerID int64, movieTitle, movieGenre string, telegramUserBoundedID *int64) error {
	if movieTitle == "" {
		return fmt.Errorf("movie title cannot be empty")
	}

	query := `
		INSERT INTO movies (telegramUsenameOwner, telegramUserOwnerID, movieTitle, movieGenre, telegramUserBoundedID)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.Exec(query, telegramUsenameOwner, telegramUserOwnerID, movieTitle, movieGenre, telegramUserBoundedID)

	if err != nil {
		return fmt.Errorf("failed to add movie: %s", err)
	}

	fmt.Println("Movie added successfully")
	return nil
}

func getMoviesHandler(telegramUsenameOwner string) (string, error) {
	query := `
		SELECT movieTitle, movieGenre
		FROM movies 
		WHERE telegramUsenameOwner = $1 
	`
	rows, err := db.Query(query, telegramUsenameOwner)
	if err != nil {
		return "", fmt.Errorf("failed to fetch movies: %s", err)
	}

	defer rows.Close()

	var result string

	for rows.Next() {
		var movieTitle, movieGenre string
		err := rows.Scan(&movieTitle, &movieGenre)

		if err != nil {
			return "", fmt.Errorf("failed to parse movies: %v", err)
		}

		if movieGenre == "" {
			result += fmt.Sprintf("Title: %s", movieTitle)
		} else {
			result += fmt.Sprintf("Title: %s, Genre: %s\n", movieTitle, movieGenre)
		}

	}

	if result == "" {
		result = "No movies found!"
	}

	return result, nil
}

func rmMovie(telegramUsenameOwner, movieTitle string) (string, error) {
	query := `
		DELETE FROM movies 
		WHERE telegramUsernameOwner = $1 
		AND movieTitle = $2
	`
	result, err := db.Exec(query, telegramUsenameOwner, movieTitle)
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
