package models

type ChatCategory struct {
	ID           int    `json:"id" bson:"categoryID"`
	CategoryName string `json:"categoryName" bson:"categoryName"`
}
