package handlers

import (
	"goozinshe/logger"
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MoviesHandler struct {
	moviesRepo   *repositories.MoviesRepository
	genresRepo   *repositories.GenresRepository
	categoryRepo *repositories.CategoryRepository
	ageRepo      *repositories.AgeRepository
	rolesRepo    *repositories.RolesRepository
}

type createMovieRequest struct {
	Title       string
	Description string
	ReleaseYear int
	Director    string
	TrailerUrl  string
	GenreIds    []int
	CategoryIds []int
	AgeIds      []int
	RoleIds     []int
}

type updateMovieRequest struct {
	Title       string
	Description string
	ReleaseYear int
	Director    string
	TrailerUrl  string
	GenreIds    []int
	CategoryIds []int
	AgeIds      []int
	RoleIds     []int
}

func NewMoviesHandler(
	moviesRepo *repositories.MoviesRepository,
	genreRepo *repositories.GenresRepository,
	categoryRepo *repositories.CategoryRepository,
	ageRepo *repositories.AgeRepository,
	rolesRepo *repositories.RolesRepository) *MoviesHandler {
	return &MoviesHandler{
		moviesRepo:   moviesRepo,
		genresRepo:   genreRepo,
		categoryRepo: categoryRepo,
		ageRepo:      ageRepo,
		rolesRepo:    rolesRepo,
	}
}

// FindById godoc
// @Summary      Find by id
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        id path int true "Movie id"
// @Success      200  {object}  models.Movie "Ok"
// @Failure      400  {object}  models.ApiError "Invalid Movie Id"
// @Router       /movies/{id} [get]
func (h *MoviesHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))

		return
	}

	movie, err := h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, movie)
}

// FindAll godoc
// @Summary      Get all movies
// @Tags         movies
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Movie "List of movies"
// @Failure      500  {object}  models.ApiError "Internal Server Error"
// @Router       /movies [get]
func (h *MoviesHandler) FindAll(c *gin.Context) {
	movies, err := h.moviesRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, movies)
}

// Create godoc
// @Summary      Create movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        title body string true "Title of the movie"
// @Param        description body string true "Description of the movie"
// @Param        releaseYear body int true "ReleaseYear of the movie"
// @Param        director body string true "Director"
// @Param        trailerUrl body string true "TrailerUrl"
// @Param      	 genreIds body []int true "Genre ids"
// @Param		 categoryIds body []int true "Category ids"
// @Param        ageIds body []int true "Age ids"
// @Success      200  {object}  object{id=int} "OK"
// @Failure      400  {object}  models.ApiError "Could not bind json"
// @Failure      500  {object}  models.ApiError
// @Router       /movies [post]
func (h *MoviesHandler) Create(c *gin.Context) {
	var request createMovieRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Could not bind json"))
		return
	}

	genres, err := h.genresRepo.FindAllByIds(c, request.GenreIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	categories, err := h.categoryRepo.FindAllByIds(c, request.CategoryIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	ages, err := h.ageRepo.FindAllByIds(c, request.AgeIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	roles, err := h.rolesRepo.FindAllByIds(c, request.RoleIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}
	movie := models.Movie{
		Title:       request.Title,
		Description: request.Description,
		ReleaseYear: request.ReleaseYear,
		Director:    request.Director,
		TrailerUrl:  request.TrailerUrl,
		Genres:      genres,
		Category:    categories,
		Ages:        ages,
		Roles:       roles,
	}

	id, err := h.moviesRepo.Create(c, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	logger := logger.GetLogger()
	logger.Info("Movie has been created", zap.Int("movie_id", id))

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// Update godoc
// @Summary      Update movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        title body string true "Title of the movie"
// @Param        description body string true "Description of the movie"
// @Param        releaseYear body int true "ReleaseYear of the movie"
// @Param        director body string true "Director"
// @Param        trailerUrl body string true "TrailerUrl"
// @Param      	 genreIds body []int true "Genre ids"
// @Param		 categoryIds body []int true "Category ids"
// @Param        ageIds body []int true "Age ids"
// @Success      200  {object}  object{id=int} "OK"
// @Failure      400  {object}  models.ApiError "Could not bind json"
// @Failure      500  {object}  models.ApiError
// @Router       /movies/{id} [put]
func (h *MoviesHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request updateMovieRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Could not bind json"))
		return
	}

	genres, err := h.genresRepo.FindAllByIds(c, request.GenreIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}
	categories, err := h.categoryRepo.FindAllByIds(c, request.CategoryIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	ages, err := h.ageRepo.FindAllByIds(c, request.AgeIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	movie := models.Movie{
		Title:       request.Title,
		Description: request.Description,
		ReleaseYear: request.ReleaseYear,
		Director:    request.Director,
		TrailerUrl:  request.TrailerUrl,
		Genres:      genres,
		Category:    categories,
		Ages:        ages,
	}

	h.moviesRepo.Update(c, id, movie)

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary      Delete movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        id path int true "Movie id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /movies/{id} [delete]
func (h *MoviesHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	h.moviesRepo.Delete(c, id)

	c.Status(http.StatusOK)
}
