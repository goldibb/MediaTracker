package services

import (
	"MediaTracker/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BookService struct {
	db *sql.DB
}
type OpenLibrarySearchResponse struct {
	NumFound int `json:"numFound"`
	Start    int `json:"start"`
	Docs     []struct {
		Title            string   `json:"title"`
		AuthorName       []string `json:"author_name,omitempty"`
		AuthorKey        []string `json:"author_key,omitempty"`
		ISBN             []string `json:"isbn,omitempty"`
		FirstPublishYear int      `json:"first_publish_year,omitempty"`
		CoverI           int      `json:"cover_i,omitempty"`
		Key              string   `json:"key"`
	} `json:"docs"`
}

func NewBookService(db *sql.DB) *BookService {
	return &BookService{db: db}
}

func (s *BookService) SearchExternalBooks(query string, page int, limit int) ([]models.Book, int, error) {
	baseUrl := "https://openlibrary.org/search.json"

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
		return nil, 0, fmt.Errorf("OpenLibrary API request failed: %w", err)
	}
	defer resp.Body.Close()

	var searchResp OpenLibrarySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, 0, fmt.Errorf("error parsing OpenLibrary response: %w", err)
	}

	totalPages := (searchResp.NumFound + limit - 1) / limit

	var books []models.Book
	for _, doc := range searchResp.Docs {
		author := ""
		if len(doc.AuthorName) > 0 {
			author = doc.AuthorName[0]
		}
		isbn := ""
		if len(doc.ISBN) > 0 {
			isbn = doc.ISBN[0]
		}

		coverURL := ""
		if doc.CoverI > 0 {
			coverURL = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", doc.CoverI)
		}
		publication_year := 0
		if doc.FirstPublishYear > 0 {
			publication_year = doc.FirstPublishYear
		}

		book := models.Book{
			Title:           doc.Title,
			Author:          author,
			PublicationYear: publication_year,
			ISBN:            isbn,
			ImageURL:        coverURL,
			Read:            false,
		}
		books = append(books, book)
	}

	return books, totalPages, nil
}

func (s *BookService) SaveBook(book models.Book) (int64, error) {
	query := `INSERT INTO books (title, author, isbn, publication_year, image_url, read) 
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id`

	var id int64
	err := s.db.QueryRow(query,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublicationYear,
		book.ImageURL,
		book.Read).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("błąd podczas zapisywania książki: %w", err)
	}

	return id, nil
}

func (s *BookService) BookExists(isbn string) (bool, error) {
	if isbn == "" {
		return false, nil
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM books WHERE isbn = $1)`
	err := s.db.QueryRow(query, isbn).Scan(&exists)

	return exists, err
}

func (s *BookService) GetBooks(search string, sort string) ([]models.Book, error) {
	query := `SELECT id, title, author, isbn, publication_year, image_url, read, created_at, updated_at FROM books WHERE 1=1`

	if search != "" {
		query += ` AND (title ILIKE $1 OR author ILIKE $1)`
	}

	// Dodajemy sortowanie
	switch sort {
	case "title_asc":
		query += ` ORDER BY title ASC`
	case "title_desc":
		query += ` ORDER BY title DESC`
	case "author_asc":
		query += ` ORDER BY author ASC`
	case "author_desc":
		query += ` ORDER BY author DESC`
	case "year_asc":
		query += ` ORDER BY publication_year ASC NULLS LAST`
	case "year_desc":
		query += ` ORDER BY publication_year DESC NULLS LAST`
	case "date_added_desc":
		query += ` ORDER BY created_at DESC`
	case "date_added_asc":
		query += ` ORDER BY created_at ASC`
	default:
		query += ` ORDER BY title ASC`
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

	var books []models.Book
	for rows.Next() {
		var book models.Book
		var createdAt, updatedAt time.Time

		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.PublicationYear,
			&book.ImageURL,
			&book.Read,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		book.CreatedAt = createdAt
		book.UpdatedAt = updatedAt
		books = append(books, book)
	}

	// Also check for errors from iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
func (s *BookService) UpdateBook(id string, book models.Book) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        UPDATE books 
        SET title = $1, author = $2, isbn = $3, publication_year = $4, image_url = $5, read = $6, updated_at = NOW()
        WHERE id = $7
    `

	result, err := tx.Exec(query,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublicationYear,
		book.ImageURL,
		book.Read,
		id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("książka o ID %s nie istnieje", id)
	}

	return tx.Commit()
}

func (s *BookService) DeleteBook(id string) error {
	query := "DELETE FROM books WHERE id = $1"
	_, err := s.db.Exec(query, id)
	return err
}
func (s *BookService) GetBookByID(id string) (models.Book, error) {
	var book models.Book
	var createdAt, updatedAt time.Time

	query := `
        SELECT id, title, author, isbn, publication_year, image_url, read, created_at, updated_at 
        FROM books 
        WHERE id = $1
    `

	err := s.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.PublicationYear,
		&book.ImageURL,
		&book.Read,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return book, fmt.Errorf("książka o ID %s nie istnieje", id)
		}
		return book, err
	}

	book.CreatedAt = createdAt
	book.UpdatedAt = updatedAt

	return book, nil
}

// Dodaj tę funkcję do BookServices.go
func (s *BookService) GetBooksByReadStatus(read bool, search string, sort string, page int, pageSize int) ([]models.Book, int, error) {
	// Najpierw policz całkowitą liczbę książek o podanym statusie
	countQuery := `SELECT COUNT(*) FROM books WHERE read = $1`

	var totalCount int
	err := s.db.QueryRow(countQuery, read).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting books: %w", err)
	}

	// Budowanie głównego zapytania z paginacją
	query := `SELECT id, title, author, isbn, publication_year, image_url, read, created_at, updated_at 
              FROM books WHERE read = $1`

	if search != "" {
		query += ` AND (title ILIKE $2 OR author ILIKE $2)`
	}

	// Dodaj sortowanie
	switch sort {
	case "title_asc":
		query += ` ORDER BY title ASC`
	case "title_desc":
		query += ` ORDER BY title DESC`
	case "author_asc":
		query += ` ORDER BY author ASC`
	case "author_desc":
		query += ` ORDER BY author DESC`
	case "year_asc":
		query += ` ORDER BY publication_year ASC NULLS LAST`
	case "year_desc":
		query += ` ORDER BY publication_year DESC NULLS LAST`
	case "date_added_desc":
		query += ` ORDER BY created_at DESC`
	case "date_added_asc":
		query += ` ORDER BY created_at ASC`
	default:
		query += ` ORDER BY title ASC`
	}

	// Dodaj paginację
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)

	var rows *sql.Rows
	if search != "" {
		rows, err = s.db.Query(query, read, "%"+search+"%")
	} else {
		rows, err = s.db.Query(query, read)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		var createdAt, updatedAt time.Time

		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.PublicationYear,
			&book.ImageURL,
			&book.Read,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		book.CreatedAt = createdAt
		book.UpdatedAt = updatedAt
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return books, totalCount, nil
}
