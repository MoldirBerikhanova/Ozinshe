package main

import (
	"context"
	"goozinshe/config"
	"goozinshe/docs"
	"goozinshe/handlers"
	"goozinshe/logger"
	"goozinshe/middlewares"
	"goozinshe/repositories"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

//
// @title           OZINSHE	API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	r := gin.Default()

	logger := logger.GetLogger()
	r.Use(
		ginzap.Ginzap(logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(logger, true),
	)

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

	moviesRepository := repositories.NewMoviesRepository(conn)
	genresRepostiroy := repositories.NewGenresRepository(conn)
	categoryRepository := repositories.NewCategoryRepository(conn)
	ageRepository := repositories.NewAgeRepository(conn)
	rolesRepository := repositories.NewRolesRepository(conn)
	moviesHandler := handlers.NewMoviesHandler(
		moviesRepository,
		genresRepostiroy,
		categoryRepository,
		ageRepository,
		rolesRepository,
	)
	usersRepository := repositories.NewUsersRepository(conn)

	genresHandler := handlers.NewGenreHanlers(genresRepostiroy)
	//imageHandlers := handlers.NewImageHandlers()
	categoryHandlers := handlers.NewCategoryHandlers(categoryRepository)
	agesHandlers := handlers.NewAgeHandler(ageRepository)
	usersHandlers := handlers.NewUsersHandlers(usersRepository)
	authHandlers := handlers.NewAuthHandlers(usersRepository)
	rolesHandlers := handlers.NewRolesHAndlers(rolesRepository)
	authorized := r.Group("")
	authorized.Use(middlewares.AuthMiddleware)

	authorized.GET("/movies/:id", moviesHandler.FindById) //http://localhost:8081/movies/:id
	authorized.GET("/movies", moviesHandler.FindAll)      //http://localhost:8081/movies/
	authorized.POST("/movies", moviesHandler.Create)
	authorized.PUT("/movies/:id", moviesHandler.Update)
	authorized.DELETE("/movies/:id", moviesHandler.Delete)

	authorized.GET("/genres/:id", genresHandler.FindById) //http://localhost:8081/genres/:id
	authorized.GET("/genres", genresHandler.FindAll)      //http://localhost:8081/genres/
	authorized.POST("/genres", genresHandler.Create)
	authorized.PUT("/genres/:id", genresHandler.Update)
	authorized.DELETE("/genres/:id", genresHandler.Delete)

	authorized.POST("/categories", categoryHandlers.Create)       //http://localhost:8081/categories/
	authorized.DELETE("/categories/:id", categoryHandlers.Delete) //http://localhost:8081/categories/:id
	authorized.GET("/categories", categoryHandlers.FindAll)
	authorized.GET("/categories/:id", categoryHandlers.FindById)
	authorized.PUT("/categories/:id", categoryHandlers.Update)

	authorized.POST("/ages", agesHandlers.HandleAddAge) //http://localhost:8081/ages/
	authorized.GET("/ages", agesHandlers.FindAll)       //http://localhost:8081/ages/:id
	authorized.GET("/ages/:id", agesHandlers.FindById)
	authorized.PUT("/ages/:id", agesHandlers.Update)
	authorized.DELETE("/ages/:id", agesHandlers.Delete)

	authorized.GET("/users", usersHandlers.FindAll)
	authorized.GET("/users/:id", usersHandlers.FindById)
	authorized.POST("/users", usersHandlers.Create)
	authorized.PUT("/users/:id", usersHandlers.Update)
	authorized.PATCH("/users/:id/changePassword", usersHandlers.ChangePassword)
	authorized.DELETE("/users/:id", usersHandlers.Delete)

	authorized.POST("/roles", rolesHandlers.Create) //http://localhost:8081/roles/
	authorized.GET("/roles", rolesHandlers.FindAll) //http://localhost:8081/roles/:id
	authorized.DELETE("/roles/:id", rolesHandlers.Delete)
	authorized.GET("/roles/:id", rolesHandlers.FindById)
	authorized.PUT("/roles/:id", rolesHandlers.Update)

	authorized.POST("/auth/signOut", authHandlers.SignOut)     //http://localhost:8081/auth/signOut
	authorized.GET("/auth/userInfo", authHandlers.GetUserInfo) //http://localhost:8081/auth/userInfo

	unauthorized := r.Group("")
	unauthorized.POST("/auth/signIn", authHandlers.SignIn) //http://localhost:8081/auth/signIn

	docs.SwaggerInfo.BasePath = "/"
	unauthorized.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler))

	logger.Info("Application starting...")
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
