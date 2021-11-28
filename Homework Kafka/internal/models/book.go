package models

type Book struct {
	ID               int    `json:"id" db:"id"`
	Title            string `json:"title" db:"title"`
	Description      string `json:"description" db:"description"`
	Genre			 string `json:"genre" db:"genre"`
	NumberOfChapters int    `json:"numberofchapters" db:"numberofchapters"`
	Author           string `json:"author" db:"author"`
	Country          string `json:"country" db:"country"`
}

type BookFilter struct {
	Query *string `json:"query"`
}
