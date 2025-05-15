package services

import (
	"MediaTracker/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ShowServices struct {
	db *sql.DB
}

func NewShowService(db *sql.DB) *ShowServices {
	return &ShowServices{db: db}
}

type TvMazeSearchResponse struct {
	Score float64     `json:"score"`
	Show  models.Show `json:"show"`
}

func (s *ShowServices) SearchExternalShows(query string, page int, limit int) ([]models.Show, int, error) {

	baseUrl := " https://api.tvmaze.com/search/shows"
	reqURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	q := reqURL.Query()
	q.Add("q", query)
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(limit))
	reqURL.RawQuery = q.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return nil, 0, fmt.Errorf("TvMaze API request failed: %w", err)
	}
	defer resp.Body.Close()

	var searchResponse []TvMazeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, 0, fmt.Errorf("failed to decode TvMaze API response: %w", err)
	}

	var shows []models.Show
	for _, doc := range searchResponse {
		shows = append(shows, doc.Show)
	}
	totalPages := (len(shows) + limit - 1) / limit
	return shows, totalPages, nil
}

func (s *ShowServices) SaveShow(show models.Show) (int64, error) {
	query := `INSERT INTO shows (name, rating_id, episodes_aired, episodes_watched, episodes_skipped, pause_status, image_url, summary, rating, genres, premiere_date, end_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	var id int64
	err := s.db.QueryRow(query,
		show.Name,
		show.RatingID,
		show.EpisodesAired,
		show.EpisodesWatched,
		show.EpisodesSkipped,
		show.PauseStatus,
		show.ImageURL,
		show.Summary,
		show.Rating,
		show.Genres,
		show.PremiereDate,
		show.EndDate,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save show: %w", err)
	}
	return id, nil
}

func (s *ShowServices) GetShows(search string, sort string) ([]models.Show, error) {
	query := `SELECT id,name, rating_id, episodes_aired, episodes_watched, episodes_skipped, pause_status, image_url, summary, rating, genres, premiere_date, end_date FROM shows WHERE 1=1`

	if search != "" {
		query += ` AND (name ILIKE $1)`
	}

	var rows *sql.Rows
	var err error
	if search != "" {
		rows, err = s.db.Query(query, "%"+search+"%")
	} else {
		rows, err = s.db.Query(query)
	}
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var shows []models.Show
	for rows.Next() {
		var show models.Show

		err = rows.Scan(
			&show.ID,
			&show.Name,
			&show.RatingID,
			&show.EpisodesAired,
			&show.EpisodesSkipped,
			&show.EpisodesWatched,
			&show.PauseStatus,
			&show.ImageURL,
			&show.Summary,
			&show.Rating,
			&show.Genres,
			&show.PremiereDate,
			&show.EndDate,
		)
		if err != nil {
			return nil, err
		}
		shows = append(shows, show)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shows, nil
}
