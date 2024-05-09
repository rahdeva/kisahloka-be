package models

import (
	"kisahloka_be/db"
	"time"
)

type Genre struct {
	GenreID   int       `json:"genre_id"`
	GenreName string    `json:"genre_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllGenres() ([]Genre, error) {
	var genres []Genre

	db := db.CreateCon()

	rows, err := db.Query("SELECT * FROM genre")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre Genre
		err := rows.Scan(&genre.GenreID, &genre.GenreName, &genre.CreatedAt, &genre.UpdatedAt)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func GetGenreDetail(genreID int) (Genre, error) {
	var genre Genre

	db := db.CreateCon()

	err := db.QueryRow("SELECT * FROM genre WHERE genre_id = ?", genreID).Scan(
		&genre.GenreID, &genre.GenreName, &genre.CreatedAt, &genre.UpdatedAt,
	)
	if err != nil {
		return Genre{}, err
	}

	return genre, nil
}

func CreateGenre(genreName string) (int64, error) {
	db := db.CreateCon()

	result, err := db.Exec("INSERT INTO genre (genre_name, created_at, updated_at) VALUES (?, ?, ?)",
		genreName, time.Now(), time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateGenre(genreID int, genreName string) (int64, error) {
	db := db.CreateCon()

	result, err := db.Exec("UPDATE genre SET genre_name = ?, updated_at = ? WHERE genre_id = ?",
		genreName, time.Now(), genreID,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func DeleteGenre(genreID int) (int64, error) {
	db := db.CreateCon()

	result, err := db.Exec("DELETE FROM genre WHERE genre_id = ?", genreID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
