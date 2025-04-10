package models

type Movie struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	RatingID           *int    `json:"rating_id"`
	Watched            bool    `json:"watched"`
	TheaterReleaseDate *string `json:"theater_release_date"`
	HomeReleaseDate    *string `json:"home_release_date"`
	ImageURL           string  `json:"image_url"`
	Description        string  `json:"description"`
}
