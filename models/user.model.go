// User model
package models

import (
	"fmt"
	"kisahloka_be/db"
	"reflect"
	"time"
)

type User struct {
	UserID    int       `json:"user_id"`
	RoleID    int       `json:"role_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllUsers retrieves all users with pagination and optional keyword search
func GetAllUsers(page, pageSize int, keyword string) (Response, error) {
	var res Response
	var arrobj reflect.Value
	var meta Meta

	con := db.CreateCon()

	// Add a WHERE clause to filter users based on the keyword
	whereClause := ""
	if keyword != "" {
		whereClause = " WHERE email LIKE '%" + keyword + "%' OR name LIKE '%" + keyword + "%'"
	}

	// Count total items in the database
	var totalItems int
	err := con.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)).Scan(&totalItems)
	if err != nil {
		return res, err
	}

	// If no items are found, return an empty response data object
	if totalItems == 0 {
		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = 0
		meta.TotalItems = totalItems

		res.Data = map[string]interface{}{
			"users": make([]interface{}, 0), // Empty slice
			"meta":  meta,
		}

		return res, nil
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Calculate the total number of pages
	totalPages := calculateTotalPages(totalItems, pageSize)

	// Check if the requested page is greater than the total number of pages
	if page > totalPages {
		return res, fmt.Errorf("requested page (%d) exceeds total number of pages (%d)", page, totalPages)
	}

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	sqlStatement := fmt.Sprintf("SELECT * FROM users %s LIMIT %d OFFSET %d", whereClause, pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj User
		err := rows.Scan(
			&obj.UserID,
			&obj.RoleID,
			&obj.Email,
			&obj.Name,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		)
		if err != nil {
			return res, err
		}

		// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
		obj.CreatedAt = obj.CreatedAt.In(loc)
		obj.UpdatedAt = obj.UpdatedAt.In(loc)

		if !arrobj.IsValid() {
			arrobj = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(obj)), 0, 0)
		}

		arrobj = reflect.Append(arrobj, reflect.ValueOf(obj))

		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = calculateTotalPages(totalItems, pageSize)
		meta.TotalItems = totalItems
	}

	res.Data = map[string]interface{}{
		"users": arrobj.Interface(),
		"meta":  meta,
	}

	return res, nil
}

// GetUserDetail retrieves details of a specific user by its ID
func GetUserDetail(userID int) (Response, error) {
	var userDetail User
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM users WHERE user_id = ?"

	row := con.QueryRow(sqlStatement, userID)

	err := row.Scan(
		&userDetail.UserID,
		&userDetail.RoleID,
		&userDetail.Email,
		&userDetail.Name,
		&userDetail.CreatedAt,
		&userDetail.UpdatedAt,
	)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
	userDetail.CreatedAt = userDetail.CreatedAt.In(loc)
	userDetail.UpdatedAt = userDetail.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"user": userDetail,
	}

	return res, nil
}

// CreateUser creates a new user with the provided data
func CreateUser(email, name string, roleID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO users (role_id, email, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	created_at := time.Now()
	updated_at := time.Now()

	result, err := stmt.Exec(
		roleID,
		email,
		name,
		created_at,
		updated_at,
	)

	if err != nil {
		return res, err
	}

	getIDLast, err := result.LastInsertId()

	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"userID":     getIDLast,
		"created_at": created_at.In(loc),
	}

	return res, nil
}

// UpdateUser updates an existing user with the provided ID and fields
func UpdateUser(userID int, updateFields map[string]interface{}) (Response, error) {
	var res Response

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Add or update the 'updated_at' field in the updateFields map
	updateFields["updated_at"] = time.Now().In(loc)
	updated_at := updateFields["updated_at"]

	con := db.CreateCon()

	// Construct the SET part of the SQL statement dynamically
	setStatement := "SET "
	values := []interface{}{}
	i := 0

	for fieldName, fieldValue := range updateFields {
		if i > 0 {
			setStatement += ", "
		}
		setStatement += fieldName + " = ?"
		values = append(values, fieldValue)
		i++
	}

	// Construct the final SQL statement
	sqlStatement := "UPDATE users " + setStatement + " WHERE user_id = ?"
	values = append(values, userID)

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(values...)
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected": rowsAffected,
		"updated_at":   updated_at,
	}

	return res, nil
}

// DeleteUser deletes a user with the provided ID
func DeleteUser(userID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM users WHERE user_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(userID)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":    rowsAffected,
		"deleted_user_id": userID,
	}

	return res, err
}
