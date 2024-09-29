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

func (h *GenreHandlers) FindAll(c *gin.Context) {
	genres := h.repo.FindAll(c)

	c.JSON(http.StatusOK, genres)
}

func (h *GenreHandlers) Create(c *gin.Context) {
	var g models.Genre
	err := c.BindJSON(&g)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	id := h.repo.Create(c, g)

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

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

	h.repo.Update(c, id, updatedGenre)

	c.Status(http.StatusOK)
}

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

	h.repo.Delete(c, id)

	c.Status(http.StatusOK)
}
