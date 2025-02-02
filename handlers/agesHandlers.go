package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

type AgeHandler struct {
	ageRepo *repositories.AgeRepository
}

func NewAgeHandler(ageRepo *repositories.AgeRepository) *AgeHandler {
	return &AgeHandler{
		ageRepo: ageRepo,
	}
}

// HandleAddAge godoc
// @Summary      HandleAddAge age
// @Tags 		 ages
// @Accept       json
// @Produce      json
// @Param request body models.Age true "Age model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid request age"
// @Failure   	 500  {object} models.ApiError
// @Router       /ages [post]
func (a *AgeHandler) HandleAddAge(c *gin.Context) {
	var age models.Age
	err := c.BindJSON(&age)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("invalid age"))
		return
	}

	id, err := a.ageRepo.Create(c, age)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// FindAll godoc
// @Summary      Get all ages
// @Tags         ages
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Age "List of age"
// @Failure      500  {object}  models.ApiError "Internal Server Error"
// @Router       /ages [get]
func (a *AgeHandler) FindAll(c *gin.Context) {
	ages, err := a.ageRepo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, ages)
}

// FindById godoc
// @Summary      Find by id
// @Tags         ages
// @Accept       json
// @Produce      json
// @Param        id path int true "Ages id"
// @Success      200  {object}  models.Age "Ok"
// @Failure      400  {object}  models.ApiError "Invalid age Id"
// @Router       /ages/{id} [get]
func (a *AgeHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("inavalid age id"))
	}

	age, err := a.ageRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, age)

}

// Update godoc
// @Summary      Update age
// @Tags 		 ages
// @Accept       json
// @Produce      json
// @Param request body models.Age true "Age model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid Age Id"
// @Failure   	 500  {object} models.ApiError
// @Router       /ages/{id} [put]
func (a *AgeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Age Id"))
		return
	}

	_, err = a.ageRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var updatedAge models.Age
	err = c.BindJSON(&updatedAge)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = a.ageRepo.Update(c, id, updatedAge)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary      Delete age
// @Tags         ages
// @Accept       json
// @Produce      json
// @Param        id path int true "Ages id"
// @Success      200  {object}  models.Age "Ok"
// @Failure      400  {object}  models.ApiError "Invalid age Id"
// @Router       /ages/{id} [delete]
func (a *AgeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid category Id"))
		return
	}

	_, err = a.ageRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = a.ageRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
