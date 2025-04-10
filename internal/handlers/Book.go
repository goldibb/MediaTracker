package handlers

import (
	"MediaTracker/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService *services.BookService
}

func CreateBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) SearchExternalBooksHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	books, err := h.bookService.SearchExternalBooks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}
func (h *BookHandler) CreateBookHandler(c *gin.Context) {
	return
}
