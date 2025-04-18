package main

var query = `
CREATE TABLE IF NOT EXISTS movies (
	id SERIAL PRIMARY KEY,
	telegramUsenameOwner VARCHAR,
	telegramUserOwnerID BIGINT NOT NULL,
	movieTitle VARCHAR NOT NULL,
	movieGenre VARCHAR,
	telegramUserBoundedID BIGINT,
	FOREIGN KEY (telegramUserOwnerID) REFERENCES users (telegramUserID) ON DELETE CASCADE);
`

var queryForUsers = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	telegramUserID BIGINT NOT NULL UNIQUE,
	telegramUsername VARCHAR);
`

var queryForSeries = `
	CREATE TABLE IF NOT EXISTS series (
		id SERIAL PRIMARY KEY,
		telegramUsernameOwner VARCHAR,
		telegramUserOwnerID BIGINT NOT NULL,
		seriesTitle VARCHAR NOT NULL,
		seriesSeasons VARCHAR NOT NULL, 
		seriesEpisodes VARCHAR NOT NULL,
		seriesGenre VARCHAR,
		telegramUserBoundedID BIGINT,
		FOREIGN KEY (telegramUserOwnerID) REFERENCES users (telegramUserID) ON DELETE CASCADE
	);
`

var addMovieQuery = `
	INSERT INTO movies (telegramUsenameOwner, telegramUserOwnerID, movieTitle, movieGenre, telegramUserBoundedID)
	VALUES ($1, $2, $3, $4, $5)
`

var addSeriesQuery = `
	INSERT INTO series (telegramusernameowner, telegramuserownerid, seriestitle, seriesseasons,seriesepisodes, seriesgenre)
	VALUES ($1, $2, $3, $4, $5, $6) 
`

var addUserQuery = `
	INSERT INTO users (telegramUserID, telegramUsername, isadmin) VALUES ($1, $2, $3)
`

var isUserExists = `
	SELECT COUNT(*) FROM users WHERE telegramUserID = $1
`

var adminQuery = `
	SELECT isadmin FROM users WHERE telegramUserID = $1
`

var userIDsSelect = `
	SELECT telegramUserID FROM users;
`

var deleteSeriesQuery = `
	DELETE FROM series
	WHERE telegramUserOwnerID = $1
	AND seriesTitle = $2
`

var getMoviesQuery = `
SELECT movieTitle, movieGenre
FROM movies 
WHERE telegramUsenameOwner = $1 
`

var getSeriesQuery = `
	SELECT seriesTitle, seriesEpisodes, seriesseasons, seriesGenre
	FROM series
	WHERE telegramUserOwnerID = $1
`

var getByGenreQuery = `
SELECT movietitle FROM movies 
WHERE telegramUsenameOwner = $1 
AND moviegenre = $2;
`

var deleteQuery = `
DELETE FROM movies 
WHERE telegramUsenameOwner ILIKE $1 
AND movietitle ILIKE $2
`
