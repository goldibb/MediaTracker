package handlers

import (
	"MediaTracker/internal/services"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

func (h *ShowHandler) SearchExternalShowsHandler(c *gin.Context) {
	query := c.PostForm("q")

	page, err := strconv.Atoi(c.PostForm("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.PostForm("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if query == "" {
		bodyBytes, _ := c.GetRawData()
		if len(bodyBytes) > 0 {
			bodyString := string(bodyBytes)
			parts := strings.Split(bodyString, "=")
			if len(parts) > 1 {
				query, _ = url.QueryUnescape(parts[1])
			}
		}
	}

	if query == "" {
		query = c.Query("q")
	}
	fmt.Printf("Search query: %s (Page: %d, Limit: %d)\n", query, page, limit)

	if query == "" {
		c.HTML(http.StatusOK, "empty_search.html", gin.H{
			"message": "Please enter a search term",
		})
		return
	}

	shows, totalPages, err := h.showService.SearchExternalShows(query, page, limit)
	if err != nil {
		fmt.Printf("Search error: %s\n", err.Error())
		c.HTML(http.StatusOK, "search_error.html", gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Printf("Found %d shows, total pages: %d\n", len(shows), totalPages)
	c.HTML(http.StatusOK, "search_results.html", gin.H{
		"shows":      shows,
		"query":      query,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}
func (h *ShowHandler) ListShowsHandler(c *gin.Context) {
	search := c.Query("search")
	sort := c.Query("sort")

	shows, err := h.showService.GetShows(search, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}
func (h *ShowHandler) GetShowsGroupedHandler(c *gin.Context) {
	sort := c.Query("sort")

	shows, err := h.showService.GetShows("", sort)
	if err != nil {
		c.HTML(http.StatusOK, "search_error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, shows)
}
