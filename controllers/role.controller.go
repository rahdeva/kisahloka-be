// Role Controller
package controllers

import (
	"kisahloka_be/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllRoles mengembalikan semua data role dengan paginasi dan pencarian kata kunci opsional
func GetAllRoles(c echo.Context) error {
	// Mendapatkan parameter query untuk paginasi
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	keyword := c.QueryParam("keyword")

	result, err := models.GetAllRoles(page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// GetRoleDetail mengembalikan detail dari suatu role berdasarkan ID-nya
func GetRoleDetail(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role_id"})
	}

	roleDetail, err := models.GetRoleDetail(roleID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, roleDetail)
}

// CreateRole membuat role baru dengan nama yang diberikan
func CreateRole(c echo.Context) error {
	var roleObj models.Role

	// Parse request body untuk mengisi struct role
	if err := c.Bind(&roleObj); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Memanggil fungsi CreateRole dari paket models
	result, err := models.CreateRole(roleObj.RoleName)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateRole memperbarui role yang ada dengan ID dan field yang diberikan
func UpdateRole(c echo.Context) error {
	// Parse request body untuk mendapatkan data update
	var updateFields map[string]interface{}
	if err := c.Bind(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Ekstrak ID dari data update
	roleID, ok := updateFields["role_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid role_id format"},
		)
	}

	// Convert id menjadi integer
	convID := int(roleID)

	// Hapus id dari map updateFields sebelum meneruskannya ke model
	delete(updateFields, "role_id")

	// Memanggil fungsi UpdateRole dari paket models
	result, err := models.UpdateRole(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteRole menghapus role dengan ID yang diberikan
func DeleteRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteRole(roleID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
