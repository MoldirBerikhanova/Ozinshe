package main

import (
	"goozinshe/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
		AllowMethods:    []string{"*"},
	}
	r.Use(cors.New(corsConfig))

	moviesHandler := handlers.NewMoviesHandler()

	r.GET("/movies/:id", moviesHandler.FindById)
	r.GET("/movies", moviesHandler.FindAll)
	r.POST("/movies", moviesHandler.Create)
	r.PUT("/movies/:id", moviesHandler.Update)
	r.DELETE("/movies/:id", moviesHandler.Delete)

	r.Run(":8081")
}
