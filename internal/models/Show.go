package models

type Show struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	RatingID        *int     `json:"rating_id"`
	EpisodesAired   int      `json:"episodes_aired"`
	EpisodesWatched []int    `json:"episodes_watched"`
	EpisodesSkipped []int    `json:"episodes_skipped"`
	PauseStatus     bool     `json:"pause_status"`
	ImageURL        string   `json:"image_url"`
	Summary         string   `json:"summary"`
	Rating          float32  `json:"rating"`
	Genres          []string `json:"genres"`
	PremiereDate    string   `json:"premiere_date"`
	EndDate         string   `json:"end_date"`
}
type PaginatedShows struct {
	Shows      []Show
	Total      int
	Page       int
	TotalPages int
	PageSize   int
}
