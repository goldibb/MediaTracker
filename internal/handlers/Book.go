package handlers

import (
	"MediaTracker/internal/models"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{
		"readBooks":       readBooks,
		"notStartedBooks": notStartedBooks,
	})
}
func (h *BookHandler) GetBookHandler(c *gin.Context) {
	id := c.Param("id")

	// Zakładając, że masz metodę GetBookByID w BookService
	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Książka nie znaleziona"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) UpdateBookHandler(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe dane książki"})
		return
	}

	// Zakładając, że masz metodę UpdateBook w BookService
	err := h.bookService.UpdateBook(id, book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie można zaktualizować książki"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Książka zaktualizowana"})
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
