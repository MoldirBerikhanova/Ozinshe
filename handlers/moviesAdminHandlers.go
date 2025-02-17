package handlers

import (
	"fmt"
	"goozinshe/logger"
	"goozinshe/models"
	"goozinshe/repositories"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MovieAdminResponseHandler struct {
	moviesAdminRepo *repositories.MoviesAdminRepository
	genresRepo      *repositories.GenresRepository
	categoryRepo    *repositories.CategoryRepository
	ageRepo         *repositories.AgeRepository
	allserieRepo    *repositories.AllSeriesRepository
}

type createMovieAdminResponseRequest struct {
	Title       string                `form:"title"`
	Description string                `form:"description"`
	ReleaseYear int                   `form:"releaseYear"`
	Director    string                `form:"director"`
	IsWatched   bool                  `form:"is_watched"`
	TrailerUrl  string                `form:"trailerUrl"`
	PosterUrl   *multipart.FileHeader `form:"posterUrl"`
	GenreIds    []int                 `form:"genreIds"`
	CategoryIds []int                 `form:"categoryIds"`
	AgeIds      []int                 `form:"ageIds"`
	AllserieIds []int                 `form:"allserieIds"`
}

type updateMovieAdminResponseRequest struct {
	Title       string                `form:"title"`
	Description string                `form:"description"`
	ReleaseYear int                   `form:"releaseYear"`
	Director    string                `form:"director"`
	IsWatched   bool                  `form:"is_watched"`
	TrailerUrl  string                `form:"trailerUrl"`
	PosterUrl   *multipart.FileHeader `form:"posterUrl"`
	GenreIds    []int                 `form:"genreIds"`
	CategoryIds []int                 `form:"categoryIds"`
	AgeIds      []int                 `form:"ageIds"`
	AllserieIds []int                 `form:"allserieIds"`
}

func NewMovieAdminResponseHandler(
	moviesAdminRepo *repositories.MoviesAdminRepository,
	genreRepo *repositories.GenresRepository,
	categoryRepo *repositories.CategoryRepository,
	ageRepo *repositories.AgeRepository,
	allserieRepo *repositories.AllSeriesRepository,
) *MovieAdminResponseHandler {
	return &MovieAdminResponseHandler{
		moviesAdminRepo: moviesAdminRepo,
		genresRepo:      genreRepo,
		categoryRepo:    categoryRepo,
		ageRepo:         ageRepo,
		allserieRepo:    allserieRepo,
	}
}

// FindById godoc
// @Summary      Find by id
// @Tags         moviesAdmin
// @Accept       json
// @Produce      json
// @Param        id path int true "Movie id"
// @Success      200  {object}  models.MovieAdminResponse
// @Failure      400  {object}  models.ApiError "Invalid Movie Id"
// @Router       /moviesAdmin/{id} [get]
func (h *MovieAdminResponseHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))

		return
	}

	movie, err := h.moviesAdminRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, movie)
}

// FindAll godoc
// @Summary      Get all moviesAdmin
// @Tags         moviesAdmin
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.MovieAdminResponse "List of movies"
// @Failure      500  {object}  models.ApiError "Internal Server Error"
// @Router       /moviesAdmin [get]
func (h *MovieAdminResponseHandler) FindAll(c *gin.Context) {
	filters := models.MovieFilters{
		SearchTerm: c.Query("search"),
		IsWatched:  c.Query("iswatched"),
		GenreId:    c.Query("genreids"),
		Sort:       c.Query("sort"),
	}
	movies, err := h.moviesAdminRepo.FindAll(c, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, movies)
}

// Create godoc
// @Summary      Create moviesAdmin
// @Tags         moviesAdmin
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
// @Router       /moviesAdmin [post]
func (h *MovieAdminResponseHandler) Create(c *gin.Context) {
	var request createMovieAdminResponseRequest

	err := c.Bind(&request)
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

	allserie, err := h.allserieRepo.FindAllByIds(c, request.AllserieIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	if request.PosterUrl == nil {
		c.JSON(http.StatusBadRequest, "Poster file is required")
		return
	}

	filename, err := h.saveMoviesPoster(c, request.PosterUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	movies := models.MovieAdminResponse{
		Title:       request.Title,
		Description: request.Description,
		ReleaseYear: request.ReleaseYear,
		Director:    request.Director,
		IsWatched:   request.IsWatched,
		TrailerUrl:  request.TrailerUrl,
		PosterUrl:   filename,
		Genres:      genres,
		Category:    categories,
		Ages:        ages,
		AllSeries:   allserie,
	}

	//IsWatched   bool        `form:"is_watched"`

	id, err := h.moviesAdminRepo.Create(c, movies)
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

func (h *MovieAdminResponseHandler) saveMoviesPoster(c *gin.Context, poster *multipart.FileHeader) (string, error) {
	filename := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(poster.Filename))
	filepath := fmt.Sprintf("images/%s", filename)
	err := c.SaveUploadedFile(poster, filepath)

	return filename, err
}

// Update godoc
// @Summary      Update moviesAdmin
// @Tags         moviesAdmin
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
// @Router       /moviesAdmin/{id} [put]
func (h *MovieAdminResponseHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	_, err = h.moviesAdminRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request updateMovieAdminResponseRequest
	err = c.Bind(&request)
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

	allseries, err := h.allserieRepo.FindAllByIds(c, request.AgeIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	filename, err := h.saveMoviesPoster(c, request.PosterUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	movie := models.MovieAdminResponse{
		Title:       request.Title,
		Description: request.Description,
		ReleaseYear: request.ReleaseYear,
		Director:    request.Director,
		IsWatched:   request.IsWatched,
		TrailerUrl:  request.TrailerUrl,
		PosterUrl:   filename,
		Genres:      genres,
		Category:    categories,
		Ages:        ages,
		AllSeries:   allseries,
	}

	err = h.moviesAdminRepo.Update(c, id, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to update movie"))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary      Delete movie
// @Tags         moviesAdmin
// @Accept       json
// @Produce      json
// @Param        id path int true "Movie id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /moviesAdmin/{id} [delete]
func (h *MovieAdminResponseHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	_, err = h.moviesAdminRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	h.moviesAdminRepo.Delete(c, id)

	c.Status(http.StatusNoContent)
}

// HandleSetWatched godoc
// @Summary      Mark moviesAdmin as watched
// @Tags         moviesAdmin
// @Produce      json
// @Param id path int true "Movie id"
// @Param isWatched query bool true "Flag value"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /moviesAdmin/{id}/setWatched [patch]
// @Security Bearer
func (h *MovieAdminResponseHandler) HandleSetWatched(c *gin.Context) {
	idStr := c.Param("movieId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	isWatchedStr := c.Query("isWatched")
	isWatched, err := strconv.ParseBool(isWatchedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid isWatched value"))
		return
	}

	err = h.moviesAdminRepo.SetWatched(c, id, isWatched)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
