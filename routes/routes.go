package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/api/v1/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Selamat Datang di KisahLoka API")
	})

	// Product
	// e.GET("/api/v1/product_variants", controllers.GetAllProductVariants)

	return e
}
