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
		books.GET("/", s.ListsBooksHandler)
		books.GET("/:id", s.GetBookHandler)
		books.POST("/", s.CreateBookHandler)
		books.PUT("/:id", s.UpdateBookHandler)
		books.DELETE("/:id", s.DeleteBookHandler)
		books.POST("/search", s.SearchExternalBooksHandler)
	}

	r.GET("/health", s.healthHandler)

	return r
}
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
