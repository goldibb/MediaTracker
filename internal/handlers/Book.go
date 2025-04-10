package handlers

import (
	"MediaTracker/internal/models"
	"MediaTracker/internal/services"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
	// Get the search query (check different possible sources)
	query := c.PostForm("q")

	// If PostForm failed, try getting from request body
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

	// Final fallback to query string
	if query == "" {
		query = c.Query("q")
	}

	// Log the query for debugging
	fmt.Printf("Search query: %s\n", query)

	if query == "" {
		c.HTML(http.StatusOK, "empty_search.html", gin.H{
			"message": "Please enter a search term",
		})
		return
	}

	books, err := h.bookService.SearchExternalBooks(query)
	if err != nil {
		fmt.Printf("Search error: %s\n", err.Error())
		c.HTML(http.StatusOK, "search_error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("Found %d books\n", len(books))
	c.HTML(http.StatusOK, "search_results.html", gin.H{
		"books": books,
	})
}
func (h *BookHandler) CreateBookHandler(c *gin.Context) {
	// Check if this is an HTMX request (form data)
	if c.GetHeader("HX-Request") != "" {
		// Process form data
		book := models.Book{
			Title:    c.PostForm("title"),
			Author:   c.PostForm("author"),
			ISBN:     c.PostForm("isbn"),
			ImageURL: c.PostForm("image_url"),
			Read:     c.PostForm("read") == "false",
		}

		// Convert publication_year from string to int
		if yearStr := c.PostForm("publication_year"); yearStr != "" {
			if year, err := strconv.Atoi(yearStr); err == nil {
				book.PublicationYear = year
			}
		}

		id, err := h.bookService.SaveBook(book)
		if err != nil {
			c.HTML(http.StatusOK, "search_error.html", gin.H{
				"error": "Failed to save book: " + err.Error(),
			})
			return
		}

		book.ID = id // Add this line

		// Return success HTML
		c.HTML(http.StatusOK, "book_added.html", gin.H{
			"book": book,
		})
	}

	// Handle normal JSON requests (for API calls)
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe dane książki: " + err.Error()})
		return
	}

	id, err := h.bookService.SaveBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie można zapisać książki: " + err.Error()})
		return
	}

	book.ID = id
	c.JSON(http.StatusCreated, book)
}
func (h *BookHandler) ListBooksHandler(c *gin.Context) {
	search := c.Query("search")
	sort := c.Query("sort")

	books, err := h.bookService.GetBooks(search, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}
func (h *BookHandler) GetBooksGroupedHandler(c *gin.Context) {
	books, err := h.bookService.GetBooks("", "title")
	if err != nil {
		c.HTML(http.StatusOK, "search_error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	readBooks := []models.Book{}
	notStartedBooks := []models.Book{}

	for _, book := range books {
		if book.Read {
			readBooks = append(readBooks, book)
		} else {
			notStartedBooks = append(notStartedBooks, book)
		}
	}

	c.HTML(http.StatusOK, "books_grouped.html", gin.H{
		"readBooks":       readBooks,
		"notStartedBooks": notStartedBooks,
	})
}
func (h *BookHandler) GetBookHandler(c *gin.Context) {
	id := c.Param("id")

	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Książka nie znaleziona"})
		return
	}

	c.JSON(http.StatusOK, book)
}
func (h *BookHandler) UpdateBookHandler(c *gin.Context) {
	id := c.Param("id")

	if c.GetHeader("HX-Request") != "" && c.ContentType() == "application/json" {
		book, err := h.bookService.GetBookByID(id)
		if err != nil {
			c.HTML(http.StatusOK, "search_error.html", gin.H{
				"error": "Book not found: " + err.Error(),
			})
			return
		}

		var partialUpdate struct {
			Read *bool `json:"Read"`
		}

		if err := c.ShouldBindJSON(&partialUpdate); err == nil && partialUpdate.Read != nil {
			book.Read = *partialUpdate.Read

			err := h.bookService.UpdateBook(id, book)
			if err != nil {
				c.HTML(http.StatusOK, "search_error.html", gin.H{
					"error": "Failed to update book: " + err.Error(),
				})
				return
			}

			if book.Read {
				c.HTML(http.StatusOK, "book_item.html", gin.H{
					"book":       book,
					"readStatus": "read",
				})
			} else {
				c.HTML(http.StatusOK, "book_item.html", gin.H{
					"book":       book,
					"readStatus": "unread",
				})
			}
			return
		}
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data: " + err.Error()})
		return
	}

	err := h.bookService.UpdateBook(id, book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update book: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated"})
}
func (h *BookHandler) DeleteBookHandler(c *gin.Context) {
	id := c.Param("id")

	// Zakładając, że masz metodę DeleteBook w BookService
	err := h.bookService.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie można usunąć książki"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Książka usunięta"})
}
