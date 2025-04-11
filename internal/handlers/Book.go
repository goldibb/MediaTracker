﻿package handlers

import (
	"MediaTracker/internal/models"
	"MediaTracker/internal/services"
	"fmt"
	"log"
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
	return NewBookHandler(bookService)
}
func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) SearchExternalBooksHandler(c *gin.Context) {
	query := c.PostForm("q")

	// Parse pagination parameters
	page, err := strconv.Atoi(c.PostForm("page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}

	limit, err := strconv.Atoi(c.PostForm("limit"))
	if err != nil || limit < 1 {
		limit = 10 // Default to 10 books per page
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

	books, totalPages, err := h.bookService.SearchExternalBooks(query, page, limit)
	if err != nil {
		fmt.Printf("Search error: %s\n", err.Error())
		c.HTML(http.StatusOK, "search_error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("Found %d books, total pages: %d\n", len(books), totalPages)
	c.HTML(http.StatusOK, "search_results.html", gin.H{
		"books":      books,
		"query":      query,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}
func (h *BookHandler) CreateBookHandler(c *gin.Context) {

	if c.GetHeader("HX-Request") != "" {
		book := models.Book{
			Title:    c.PostForm("title"),
			Author:   c.PostForm("author"),
			ISBN:     c.PostForm("isbn"),
			ImageURL: c.PostForm("image_url"),
			Read:     c.PostForm("read") == "false",
		}

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

		book.ID = id

		c.HTML(http.StatusOK, "book_added.html", gin.H{
			"book": book,
		})
		return
	}
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {

		log.Printf("Invalid JSON book data: %v", err)

		if c.GetHeader("HX-Request") != "" {
			c.HTML(http.StatusBadRequest, "search_error.html", gin.H{
				"error": "Please provide valid book data",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe dane książki: " + err.Error()})
		}
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

	sort := c.Query("sort")
	if sortCookie, err := c.Cookie("book_sort"); err == nil && sort == "" {
		sort = sortCookie
	}
	if sort == "" {
		sort = "title_asc"
	}

	c.SetCookie("book_sort", sort, 86400, "/", "", false, true)

	books, err := h.bookService.GetBooks("", sort)
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
		"currentSort":     sort,
	})
}
func (h *BookHandler) UpdateBookHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("UpdateBookHandler called with ID: %s, content-type: %s\n", id, c.ContentType())

	if c.GetHeader("HX-Request") != "" {
		book, err := h.bookService.GetBookByID(id)
		if err != nil {
			c.HTML(http.StatusOK, "search_error.html", gin.H{
				"error": "Book not found: " + err.Error(),
			})
			return
		}

		var readValue bool

		if c.ContentType() == "application/json" {

			var partialUpdate struct {
				Read *bool `json:"Read"`
			}
			if err := c.ShouldBindJSON(&partialUpdate); err == nil && partialUpdate.Read != nil {
				readValue = *partialUpdate.Read
			} else {
				c.HTML(http.StatusOK, "search_error.html", gin.H{
					"error": "Invalid JSON data: " + err.Error(),
				})
				return
			}
		} else {

			readStr := c.PostForm("Read")
			readValue = readStr == "true"
		}

		sort := c.Query("sort")
		if sortCookie, err := c.Cookie("book_sort"); err == nil && sort == "" {
			sort = sortCookie
		}
		if sort == "" {
			sort = "title_asc"
		}
		book.Read = readValue
		err = h.bookService.UpdateBook(id, book)
		if err != nil {
			c.HTML(http.StatusOK, "search_error.html", gin.H{
				"error": "Failed to update book: " + err.Error(),
			})
			return
		}

		fmt.Printf("Book status updated - ID: %s, Title: %s, Read: %t\n", id, book.Title, book.Read)

		books, err := h.bookService.GetBooks("", "sort")
		if err != nil {
			c.HTML(http.StatusOK, "search_error.html", gin.H{
				"error": "Failed to retrieve books: " + err.Error(),
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
			"currentSort":     sort,
		})
		return
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

	err := h.bookService.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie można usunąć książki"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Książka usunięta"})
}
