package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandlers struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryHandlers(categoryRepo *repositories.CategoryRepository) *CategoryHandlers {
	return &CategoryHandlers{
		categoryRepo: categoryRepo,
	}
}

// Create godoc
// @Summary      Create category
// @Tags 		 categories
// @Accept       json
// @Produce      json
// @Param request body models.Category true "Category model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid request category"
// @Failure   	 500  {object} models.ApiError
// @Router       /categories [post]
func (h *CategoryHandlers) Create(c *gin.Context) {
	var category models.Category
	err := c.BindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request category")
		return
	}

	id, err := h.categoryRepo.Create(c, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// FindById godoc
// @Summary      Find by id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id path int true "Category id"
// @Success      200  {object}  models.Category "Ok"
// @Failure      400  {object}  models.ApiError "Invalid Movie Id"
// @Router       /categories/{id} [get]
func (h *CategoryHandlers) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid  categoryId"))
		return
	}

	_, err = h.categoryRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

}

// FindAll godoc
// @Summary      Get all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Movie "List of movies"
// @Failure      500  {object}  models.ApiError "Internal Server Error"
// @Router       /categories [get]
func (h *CategoryHandlers) FindAll(c *gin.Context) {
	categories, err := h.categoryRepo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Update godoc
// @Summary      Update category
// @Tags 		 categories
// @Accept       json
// @Produce      json
// @Param request body models.Category true "Category model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid Category Id"
// @Failure   	 500  {object} models.ApiError
// @Router       /categories/{id} [put]
func (h *CategoryHandlers) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
		return
	}

	_, err = h.categoryRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var updatedCategory models.Category
	err = c.BindJSON(&updatedCategory)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = h.categoryRepo.Update(c, id, updatedCategory)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary      Delete category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id path int true "Category id"
// @Success      200  {object}  models.Category "Ok"
// @Failure      400  {object}  models.ApiError "Invalid category Id"
// @Router       /categories/{id} [delete]
func (h *CategoryHandlers) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid category Id"))
		return
	}

	_, err = h.categoryRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = h.categoryRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
