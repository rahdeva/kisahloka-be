package controllers

import (
	"kisahloka_be/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetHomeData(c echo.Context) error {
	genres, err := models.GetHomeData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, genres)
}
