package models

type Book struct {
	Id        string  `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
	Date      string  `json:"date"`
}
