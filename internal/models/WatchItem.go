package models

type WatchItem struct {
	ID          int
	Name        string
	RatingID    *int
	Href        string
	ImageURL    string
	Description string
}
