package main

import (
	"context"
	"goozinshe/config"
	"goozinshe/handlers"
	"goozinshe/repositories"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
		AllowMethods:    []string{"*"},
	}
	r.Use(cors.New(corsConfig))

	err := loadConfig()
	if err != nil {
		panic(err)
	}

	conn, err := connectToDb()
	if err != nil {
		panic(err)
	}

	moviesRepository := repositories.NewMoviesRepository()
	genresRepostiroy := repositories.NewGenresRepository(conn)
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

	r.Run(config.Config.AppHost)
}

func loadConfig() error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var mapConfig config.MapConfig
	err = viper.Unmarshal(&mapConfig)
	if err != nil {
		return err
	}

	config.Config = &mapConfig

	return nil
}

func connectToDb() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
