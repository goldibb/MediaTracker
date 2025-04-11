package server

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5137"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-HX-Request",
			"X-HX-Current-URL",
			"X-HX-Target",
		},
		AllowCredentials: true,
	}))

	books := r.Group("/api/books")
	{
		books.GET("/", s.bookHandler.ListBooksHandler)
		books.GET("/grouped", s.bookHandler.GetBooksGroupedHandler)
		books.GET("/:id", s.bookHandler.GetBooksGroupedHandler)
		books.POST("/", s.bookHandler.CreateBookHandler)
		books.PUT("/:id", s.bookHandler.UpdateBookHandler)
		r.DELETE("/api/books/:id", s.bookHandler.DeleteBookHandler)
		books.POST("/search", s.bookHandler.SearchExternalBooksHandler)
	}
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})
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
