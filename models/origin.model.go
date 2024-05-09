// Origin model
package models

import (
	"fmt"
	"kisahloka_be/db"
	"reflect"
	"time"
)

type Origin struct {
	OriginID   int       `json:"origin_id"`
	OriginName string    `json:"origin_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetAllOrigins retrieves all origins with pagination and optional keyword filtering
func GetAllOrigins(page, pageSize int, keyword string) (Response, error) {
	var res Response
	var arrobj reflect.Value
	var meta Meta

	con := db.CreateCon()

	// Add a WHERE clause to filter origins based on the keyword
	whereClause := ""
	if keyword != "" {
		whereClause = " WHERE origin_name LIKE '%" + keyword + "%'"
	}

	// Count total items in the database
	var totalItems int
	err := con.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM origin %s", whereClause)).Scan(&totalItems)
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
			"origins": make([]interface{}, 0), // Empty slice
			"meta":    meta,
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
	sqlStatement := fmt.Sprintf("SELECT * FROM origin %s LIMIT %d OFFSET %d", whereClause, pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj Origin
		err := rows.Scan(
			&obj.OriginID,
			&obj.OriginName,
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
		"origins": arrobj.Interface(),
		"meta":    meta,
	}

	return res, nil
}

// GetOriginDetail retrieves details of a specific origin by ID
func GetOriginDetail(originID int) (Response, error) {
	var originDetail Origin
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM origin WHERE origin_id = ?"

	row := con.QueryRow(sqlStatement, originID)

	err := row.Scan(
		&originDetail.OriginID,
		&originDetail.OriginName,
		&originDetail.CreatedAt,
		&originDetail.UpdatedAt,
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
	originDetail.CreatedAt = originDetail.CreatedAt.In(loc)
	originDetail.UpdatedAt = originDetail.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"origin": originDetail,
	}

	return res, nil
}

// CreateOrigin creates a new origin
func CreateOrigin(originName string) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO origin (origin_name, created_at, updated_at) VALUES (?, ?, ?)"

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
		originName,
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
		"getIDLast":  getIDLast,
		"created_at": created_at.In(loc),
	}

	return res, nil
}

// UpdateOrigin updates an existing origin
func UpdateOrigin(originID int, updateFields map[string]interface{}) (Response, error) {
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
	sqlStatement := "UPDATE origin " + setStatement + " WHERE origin_id = ?"
	values = append(values, originID)

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

// DeleteOrigin deletes an origin by ID
func DeleteOrigin(originID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM origin WHERE origin_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(originID)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":      rowsAffected,
		"deleted_origin_id": originID,
	}

	return res, err
}
