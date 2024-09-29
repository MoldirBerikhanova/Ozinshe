package main

import (
	"goozinshe/handlers"
	"goozinshe/repositories"

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

	moviesRepository := repositories.NewMoviesRepository()
	genresRepostiroy := repositories.NewGenresRepository()
	moviesHandler := handlers.NewMoviesHandler(
		moviesRepository,
		genresRepostiroy,
	)
	genresHandler := handlers.NewGenreHanlers(genresRepostiroy)

	r.GET("/movies/:id", moviesHandler.FindById)
	r.GET("/movies", moviesHandler.FindAll)
	r.POST("/movies", moviesHandler.Create)
	r.PUT("/movies/:id", moviesHandler.Update)
	r.DELETE("/movies/:id", moviesHandler.Delete)

	r.GET("/genres/:id", genresHandler.FindById)
	r.GET("/genres", genresHandler.FindAll)
	r.POST("/genres", genresHandler.Create)
	r.PUT("/genres/:id", genresHandler.Update)
	r.DELETE("/genres/:id", genresHandler.Delete)

	r.Run(":8081")
}
