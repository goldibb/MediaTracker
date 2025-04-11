package server

import (
	"MediaTracker/internal/database"
	"MediaTracker/internal/handlers"
	"MediaTracker/internal/services"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port        int
	bookHandler *handlers.BookHandler
	db          database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	bookService := services.NewBookService(db.GetDB())
	bookHandler := handlers.NewBookHandler(bookService)
	NewServer := &Server{
		port:        port,
		bookHandler: bookHandler,
		db:          db,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
