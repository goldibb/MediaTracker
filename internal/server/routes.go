package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	books := r.Group("/api/books")
	{
		books.GET("/", s.bookHandler.ListBooksHandler)              // Zmieniono z GetBooksHandler
		books.GET("/grouped", s.bookHandler.GetBooksGroupedHandler) // Dodano nowy endpoint
		books.GET("/:id", s.bookHandler.GetBookHandler)             // Ten handler musisz zaimplementować
		books.POST("/", s.bookHandler.CreateBookHandler)
		books.PUT("/:id", s.bookHandler.UpdateBookHandler)    // Ten handler musisz zaimplementować
		books.DELETE("/:id", s.bookHandler.DeleteBookHandler) // Ten handler musisz zaimplementować
		books.POST("/search", s.bookHandler.SearchExternalBooksHandler)
	}

	r.Static("/static", "./frontend")
	r.LoadHTMLGlob("frontend/*.html")

	r.GET("/books", func(c *gin.Context) {
		c.HTML(http.StatusOK, "readlist.html", nil)
	})

	r.GET("/books/add", func(c *gin.Context) {
		c.HTML(http.StatusOK, "BookAdd.html", nil)
	})

	r.GET("/health", s.healthHandler)

	return r
}
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
