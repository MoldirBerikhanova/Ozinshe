package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandlers struct {
	userRepo *repositories.UsersRepository
}

func NewUsersHandlers(userRepo *repositories.UsersRepository) *UsersHandlers {
	return &UsersHandlers{userRepo: userRepo}
}

type createUserRequest struct {
	Name        string
	Email       string
	Password    string
	PhoneNumber *int
	Birthday    *time.Time
}

type updateUserRequest struct {
	Name        string
	Email       string
	PhoneNumber *int
	Birthday    *time.Time
}

type changePasswordRequest struct {
	Password string
}

type userResponse struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	PhoneNumber *int       `json:"phonenumber"`
	Birthday    *time.Time `json:"birthday"`
}

// FindById godoc
// @Tags users
// @Summary      Find users by id
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Success      200  {array} handlers.userResponse "OK"
// @Failure   	 400  {object} models.ApiError "Invalid user id"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [get]
func (h *UsersHandlers) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user Id"))
		return
	}

	user, err := h.userRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	r := userResponse{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Birthday:    user.Birthday,
	}

	c.JSON(http.StatusOK, r)
}

// FindAll godoc
// @Tags users
// @Summary      Get users list
// @Accept       json
// @Produce      json
// @Success      200  {array} handlers.userResponse "OK"
// @Failure   	 500  {object} models.ApiError
// @Router       /users [get]
func (h *UsersHandlers) FindAll(c *gin.Context) {
	users, err := h.userRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not load users"))
		return
	}

	dtos := make([]userResponse, 0, len(users))
	for _, u := range users {
		r := userResponse{
			Id:          u.Id,
			Name:        u.Name,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			Birthday:    u.Birthday,
		}
		dtos = append(dtos, r)
	}

	c.JSON(http.StatusOK, dtos)
}

// Create godoc
// @Tags users
// @Summary      Create user
// @Accept       json
// @Produce      json
// @Param request body handlers.createUserRequest true "User data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /users [post]
func (h *UsersHandlers) Create(c *gin.Context) {
	var request createUserRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid payload"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to hash password"))
		return
	}

	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		PhoneNumber:  request.PhoneNumber,
		Birthday:     request.Birthday,
	}

	id, err := h.userRepo.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not create user"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Update godoc
// @Tags users
// @Summary      Update user
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Param request body handlers.updateUserRequest true "User data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [put]
func (h *UsersHandlers) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user Id"))
		return
	}

	var request updateUserRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	user, err := h.userRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	user.Name = request.Name
	user.Email = request.Email
	user.PhoneNumber = request.PhoneNumber
	user.Birthday = request.Birthday

	err = h.userRepo.Update(c, id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Tags users
// @Summary      Delete user
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [delete]
func (h *UsersHandlers) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user Id"))
		return
	}

	_, err = h.userRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	err = h.userRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// ChangePassword godoc
// @Tags users
// @Summary      Change user password
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Param request body handlers.changePasswordRequest true "Password data"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id}/changePassword [patch]
func (h *UsersHandlers) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user Id"))
		return
	}

	var request changePasswordRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to hash password"))
		return
	}

	user, err := h.userRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	user.PasswordHash = string(passwordHash)

	err = h.userRepo.Update(c, id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
