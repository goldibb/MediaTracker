package models

type Rating struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}
