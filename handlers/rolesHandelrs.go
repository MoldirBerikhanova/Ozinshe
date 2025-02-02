package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RolesHandlers struct {
	rolesRepo *repositories.RolesRepository
}

func NewRolesHAndlers(rolesRepo *repositories.RolesRepository) *RolesHandlers {
	return &RolesHandlers{
		rolesRepo: rolesRepo,
	}
}

// Create godoc
// @Summary      Create role
// @Tags 		 roles
// @Accept       json
// @Produce      json
// @Param request body models.Roles true "Roles model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid request roles"
// @Failure   	 500  {object} models.ApiError
// @Router       /roles [post]
func (h *RolesHandlers) Create(c *gin.Context) {
	var role models.Roles
	err := c.BindJSON(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("roles couldnt create"))
		return
	}

	id, err := h.rolesRepo.Create(c, role)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// FindAll godoc
// @Summary      Get all roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Roles "List of roles"
// @Failure      500  {object}  models.ApiError "Internal Server Error"
// @Router       /roles [get]
func (h *RolesHandlers) FindAll(c *gin.Context) {
	roles, err := h.rolesRepo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// FindById godoc
// @Summary      Find by id
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id path int true "Roles id"
// @Success      200  {object}  models.Roles "Ok"
// @Failure      400  {object}  models.ApiError "Invalid Role Id"
// @Router       /roles/{id} [get]
func (h *RolesHandlers) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid  rolesId"))
		return
	}

	role, err := h.rolesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, role)
}

// Update godoc
// @Summary      Update roles
// @Tags 		 roles
// @Accept       json
// @Produce      json
// @Param request body models.Roles true "Roles model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid Role Id"
// @Failure   	 500  {object} models.ApiError
// @Router       /roles/{id} [put]
func (h *RolesHandlers) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("invalid request"))
		return
	}

	_, err = h.rolesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("role not found"))
		return
	}
	var updatedroles models.Roles
	err = c.BindJSON(&updatedroles)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = h.rolesRepo.Update(c, id, updatedroles)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)

}

// Delete godoc
// @Summary      Delete roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id path int true "Roles id"
// @Success      200  {object}  models.Roles "Ok"
// @Failure      400  {object}  models.ApiError "Invalid role Id"
// @Router       /roles/{id} [delete]
func (h *RolesHandlers) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("invalid request"))
		return
	}

	_, err = h.rolesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("role not found"))
		return
	}

	err = h.rolesRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
