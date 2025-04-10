package models

import "time"

type Book struct {
	ID              int64
	Title           string
	Author          string
	ISBN            string
	PublicationYear int
	Read            bool
	ImageURL        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PaginatedBooks struct {
	Books      []Book
	Total      int
	Page       int
	TotalPages int
	PageSize   int
}
