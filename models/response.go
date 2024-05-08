package models

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

type Meta struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalPages int `json:"total_page"`
	TotalItems int `json:"total_items"`
}

func ResponseData(typeName string) interface{} {
	switch typeName {
	case "type":
		return &Type{}
	default:
		return nil
	}
}

func calculateTotalPages(totalItems, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	totalPages := totalItems / pageSize
	if totalItems%pageSize > 0 {
		totalPages++
	}
	return totalPages
}
