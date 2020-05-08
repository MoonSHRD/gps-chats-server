package models

type ChatCategory struct {
	Id           int    `json:"id" db:"id"`
	CategoryName string `json:"categoryName" db:"categoryname"`
}
