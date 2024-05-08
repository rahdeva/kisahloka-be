package models

import (
	"fmt"
	"kisahloka_be/db"
	"reflect"
	"time"
)

type Type struct {
	TypeID    int       `json:"type_id"`
	TypeName  string    `json:"type_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllTypes(page, pageSize int, keyword string) (Response, error) {
	var res Response
	var arrobj reflect.Value
	var meta Meta

	con := db.CreateCon()

	// Add a WHERE clause to filter types based on the keyword
	whereClause := ""
	if keyword != "" {
		whereClause = " WHERE type_name LIKE '%" + keyword + "%'"
	}

	// Count total items in the database
	var totalItems int
	err := con.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM type %s", whereClause)).Scan(&totalItems)
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
			"types": make([]interface{}, 0), // Empty slice
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
	sqlStatement := fmt.Sprintf("SELECT * FROM type %s LIMIT %d OFFSET %d", whereClause, pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj Type
		err := rows.Scan(
			&obj.TypeID,
			&obj.TypeName,
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
		"types": arrobj.Interface(),
		"meta":  meta,
	}

	return res, nil
}

func GetTypeDetail(typeID int) (Response, error) {
	var typeDetail Type
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM type WHERE type_id = ?"

	row := con.QueryRow(sqlStatement, typeID)

	err := row.Scan(
		&typeDetail.TypeID,
		&typeDetail.TypeName,
		&typeDetail.CreatedAt,
		&typeDetail.UpdatedAt,
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
	typeDetail.CreatedAt = typeDetail.CreatedAt.In(loc)
	typeDetail.UpdatedAt = typeDetail.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"type": typeDetail,
	}

	return res, nil
}

func CreateType(typeName string) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO type (type_name, created_at, updated_at) VALUES (?, ?, ?)"

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
		typeName,
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

func UpdateType(typeID int, typeName string) (Response, error) {
	var res Response

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Add or update the 'updated_at' field
	updated_at := time.Now().In(loc)

	con := db.CreateCon()

	sqlStatement := "UPDATE type SET type_name = ?, updated_at = ? WHERE type_id = ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(typeName, updated_at, typeID)
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

func DeleteType(typeID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM type WHERE type_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(typeID)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":    rowsAffected,
		"deleted_type_id": typeID,
	}

	return res, err
}
