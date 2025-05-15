package handlers

import (
	"MediaTracker/internal/services"
)

type ShowHandler struct {
	showService *services.ShowServices
}

func CreateShowHandler(showService *services.ShowServices) *ShowHandler {
	return NewShowHandler(showService)
}
func NewShowHandler(showService *services.ShowServices) *ShowHandler {
	return &ShowHandler{
		showService: showService,
	}
}
