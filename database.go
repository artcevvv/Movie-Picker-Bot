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

func getMoviesHandler(telegramUsenameOwner string) ([]map[string]string, error) {
	query := `
		SELECT movieTitle, movieGenre
		FROM movies 
		WHERE telegramUsenameOwner = $1 
	`
	rows, err := db.Query(query, telegramUsenameOwner)
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
	query := `
		SELECT movietitle FROM movies 
		WHERE telegramUsenameOwner = $1 
		AND moviegenre = $2;
	`

	rows, err := db.Query(query, telegramUsenameOwner, movieGenre)
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

// func getGenres(telegramUsenameOwner string) ([]map[string]string, error) {
// 	query := `
// 		SELECT moviegenre FROM movies
// 		WHERE telegramUsenameOwner = $1
// 	`

// 	rows, err := db.Query(query, telegramUsenameOwner)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch genres: %v", err)
// 	}

// 	defer rows.Close()

// 	var genres []map[string]string

// 	for rows.Next() {
// 		var movieGenre string
// 		err := rows.Scan(&movieGenre)

// 		if err != nil {
// 			return nil, fmt.Errorf("failed to parse genres: %v")
// 		}
// 	}
// }

func rmMovie(telegramUsenameOwner, movieTitle string) (string, error) {
	query := `
		DELETE FROM movies 
		WHERE telegramUsenameOwner ILIKE $1 
		AND movietitle ILIKE $2
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
