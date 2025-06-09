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
	// Books handlers
	books := r.Group("/api/books")
	{
		books.GET("/", s.bookHandler.ListBooksHandler)
		books.GET("/grouped", s.bookHandler.GetBooksGroupedHandler)
		books.GET("/:id", s.bookHandler.GetBooksGroupedHandler)
		books.POST("/", s.bookHandler.CreateBookHandler)
		books.PUT("/:id", s.bookHandler.UpdateBookHandler)
		books.DELETE("/:id", s.bookHandler.DeleteBookHandler)
		books.POST("/search", s.bookHandler.SearchExternalBooksHandler)
		books.GET("/edit/:id", s.bookHandler.EditBookHandler)
		books.POST("/edit/:id", s.bookHandler.UpdateBookDetailsHandler)
	}

	// Shows handlers
	shows := r.Group("/api/shows")
	{
		shows.GET("/", s.showHandler.ListShowsHandler)
		shows.GET("/grouped")
		shows.GET("/:id")
		shows.POST("/")
		shows.PUT("/:id")
		shows.DELETE("/:id")
		shows.POST("/search")
		shows.GET("/edit/:id")
		shows.POST("/edit/:id")
	}

	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})
	r.Static("/static", "./frontend")
	r.LoadHTMLGlob("frontend/*.html")

	// Views
	views := r.Group("/books")
	{
		views.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "readlist.html", nil)
		})

		views.GET("/add", func(c *gin.Context) {
			c.HTML(http.StatusOK, "BookAdd.html", nil)
		})
	}
	showViews := r.Group("/shows")
	{
		showViews.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "showlist.html", nil)
		})
	}
	r.GET("/health", s.healthHandler)

	return r
}
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
