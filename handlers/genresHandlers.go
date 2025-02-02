package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GenreHandlers struct {
	repo *repositories.GenresRepository
}

func NewGenreHanlers(repo *repositories.GenresRepository) *GenreHandlers {
	return &GenreHandlers{
		repo: repo,
	}
}

// FindById godoc
// @Summary      Find by id
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id path int true "Genre id"
// @Success      200  {object}  models.Genre "Ok"
// @Failure      400  {object}  models.ApiError "Invalid Movie Id"
// @Router       /genres/{id} [get]
func (h *GenreHandlers) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
		return
	}

	genre, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, genre)
}

// FindAll godoc
// @Summary      Get all genres
// @Tags         genres
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Genre "List of genres"
// @Failure      500  {object}  models.ApiError "Internal Server Error"
// @Router       /genres [get]
func (h *GenreHandlers) FindAll(c *gin.Context) {
	genres, err := h.repo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, genres)
}

// Create godoc
// @Summary      Create genre
// @Tags 		 genres
// @Accept       json
// @Produce      json
// @Param request body models.Genre true "Genre model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid request category"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres [post]
func (h *GenreHandlers) Create(c *gin.Context) {
	var g models.Genre
	err := c.BindJSON(&g)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	id, err := h.repo.Create(c, g)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// Update godoc
// @Summary      Update genre
// @Tags 		 genres
// @Accept       json
// @Produce      json
// @Param request body models.Genre true "Genre model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid Genre Id"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres/{id} [put]
func (h *GenreHandlers) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var updatedGenre models.Genre
	err = c.BindJSON(&updatedGenre)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = h.repo.Update(c, id, updatedGenre)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary      Delete genre
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id path int true "Genre id"
// @Success      200  {object}  models.Genre "Ok"
// @Failure      400  {object}  models.ApiError "Invalid genre Id"
// @Router       /genres/{id} [delete]
func (h *GenreHandlers) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = h.repo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
