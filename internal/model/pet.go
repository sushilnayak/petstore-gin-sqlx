package model

type Pet struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Status     string `json:"status" db:"status"`
	CategoryID int64  `json:"categoryId" db:"category_id"`
}
