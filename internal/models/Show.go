package models

type Show struct {
	ID              int      `json:"id" db:"id"`
	Name            string   `json:"name" db:"name"`
	RatingID        *int     `json:"rating_id" db:"rating_id"`
	EpisodesAired   int      `json:"episodes_aired" db:"episodes_aired"`
	EpisodesWatched []int    `json:"episodes_watched" db:"episodes_watched"`
	EpisodesSkipped []int    `json:"episodes_skipped" db:"episodes_skipped"`
	PauseStatus     bool     `json:"pause_status" db:"pause_status"`
	ImageURL        string   `json:"image_url" db:"image_url"`
	Summary         string   `json:"summary" db:"summary"`
	Rating          float32  `json:"rating" db:"rating"`
	Genres          []string `json:"genres" db:"genres"`
	PremiereDate    string   `json:"premiere_date" db:"premiere_date"`
	EndDate         string   `json:"end_date" db:"end_date"`
}
type PaginatedShows struct {
	Shows      []Show
	Total      int
	Page       int
	TotalPages int
	PageSize   int
}
