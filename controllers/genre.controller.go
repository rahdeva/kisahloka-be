package controllers

import (
	"kisahloka_be/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllGenres(c echo.Context) error {
	genres, err := models.GetAllGenres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, genres)
}

func GetGenreDetail(c echo.Context) error {
	genreID, err := strconv.Atoi(c.Param("genre_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid genre_id"})
	}

	genre, err := models.GetGenreDetail(genreID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, genre)
}

func CreateGenre(c echo.Context) error {
	var genre struct {
		GenreName string `json:"genre_name"`
	}

	if err := c.Bind(&genre); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	id, err := models.CreateGenre(genre.GenreName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int64{"genre_id": id})
}

func UpdateGenre(c echo.Context) error {
	genreID, err := strconv.Atoi(c.Param("genre_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid genre_id"})
	}

	var genre struct {
		GenreName string `json:"genre_name"`
	}

	if err := c.Bind(&genre); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	rowsAffected, err := models.UpdateGenre(genreID, genre.GenreName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int64{"rows_affected": rowsAffected})
}

func DeleteGenre(c echo.Context) error {
	genreID, err := strconv.Atoi(c.Param("genre_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid genre_id"})
	}

	rowsAffected, err := models.DeleteGenre(genreID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int64{"rows_affected": rowsAffected})
}
