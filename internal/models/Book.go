package models

import "time"

type Book struct {
	ID              int64     `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Author          string    `json:"author" db:"author"`
	ISBN            string    `json:"isbn" db:"isbn"`
	PublicationYear int       `json:"publication_year" db:"publication_year"`
	Read            bool      `json:"read" db:"read"`
	ImageURL        string    `json:"image_url" db:"image_url"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type PaginatedBooks struct {
	Books      []Book
	Total      int
	Page       int
	TotalPages int
	PageSize   int
}
