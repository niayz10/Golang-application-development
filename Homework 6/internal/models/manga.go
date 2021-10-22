package models

type Manga struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Genre			 string `json:"genre"`
	NumberOfChapters int    `json:"numofchapters"`
	Author           string `json:"author"`
	Country          string `json:"country"`
}
