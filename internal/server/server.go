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
	showHandler *handlers.ShowHandler
	db          database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	bookService := services.NewBookService(db.GetDB())
	bookHandler := handlers.NewBookHandler(bookService)
	showService := services.NewShowService(db.GetDB())
	showHandler := handlers.NewShowHandler(showService)
	NewServer := &Server{
		port:        port,
		bookHandler: bookHandler,
		showHandler: showHandler,
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
