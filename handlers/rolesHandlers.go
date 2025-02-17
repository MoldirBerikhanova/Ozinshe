package handlers

import (
	"fmt"
	"goozinshe/models"
	"goozinshe/repositories"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RolesHandlers struct {
	rolesRepo       *repositories.RolesRepository
	userRepo        *repositories.UsersRepository
	moviesAdminRepo *repositories.MoviesAdminRepository
	genresRepo      *repositories.GenresRepository
	categoryRepo    *repositories.CategoryRepository
	ageRepo         *repositories.AgeRepository
	allserieRepo    *repositories.AllSeriesRepository
}

func NewRolesHandlers(
	rolesRepo *repositories.RolesRepository,
	userRepo *repositories.UsersRepository,
	moviesAdminRepo *repositories.MoviesAdminRepository,
	genreRepo *repositories.GenresRepository,
	categoryRepo *repositories.CategoryRepository,
	ageRepo *repositories.AgeRepository,
	allserieRepo *repositories.AllSeriesRepository) *RolesHandlers {
	return &RolesHandlers{
		rolesRepo:       rolesRepo,
		userRepo:        userRepo,
		moviesAdminRepo: moviesAdminRepo,
		genresRepo:      genreRepo,
		categoryRepo:    categoryRepo,
		ageRepo:         ageRepo,
		allserieRepo:    allserieRepo}
}

type createRolesRequest struct {
	Name        string
	Email       string
	Password    string
	PhoneNumber *int
	Birthday    *time.Time
}

type updateRolesRequest struct {
	Name        string
	Email       string
	PhoneNumber *int
	Birthday    *time.Time
}

type changeRolesPasswordRequest struct {
	Password string
}

type rolesResponse struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	PhoneNumber *int       `json:"phonenumber"`
	Birthday    *time.Time `json:"birthday"`
}

// FindById godoc
// @Tags         roles - это выполняет роль Админа
// @Summary      Find role by id
// @Accept       json
// @Produce      json
// @Param id path int true "Roles id"
// @Success      200  {array}  handlers.rolesResponse  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid role id"
// @Failure   	 404  {object} models.ApiError "Role not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /roles/{id} [get]
func (h *RolesHandlers) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid role Id"))
		return
	}

	role, err := h.rolesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("Role not found"))
		return
	}

	r := rolesResponse{
		Id:          role.Id,
		Name:        role.Name,
		Email:       role.Email,
		PhoneNumber: role.PhoneNumber,
		Birthday:    role.Birthday,
	}

	c.JSON(http.StatusOK, r)
}

// FindAll godoc
// @Tags         roles - это выполняет роль Админа
// @Summary      Get roles list
// @Accept       json
// @Produce      json
// @Success      200  {array} handlers.rolesResponse "OK"
// @Failure   	 500  {object} models.ApiError
// @Router       /roles [get]
func (h *RolesHandlers) FindAll(c *gin.Context) {
	roles, err := h.rolesRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not load roles"))
		return
	}

	dtos := make([]rolesResponse, 0, len(roles))
	for _, r := range roles {
		r := rolesResponse{
			Id:          r.Id,
			Name:        r.Name,
			Email:       r.Email,
			PhoneNumber: r.PhoneNumber,
			Birthday:    r.Birthday,
		}
		dtos = append(dtos, r)
	}

	c.JSON(http.StatusOK, dtos)
}

// Create godoc
// @Tags         roles - это выполняет роль Админа
// @Summary      Create role
// @Accept       json
// @Produce      json
// @Param request body handlers.createRolesRequest true "Roles data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /roles [post]
func (h *RolesHandlers) Create(c *gin.Context) {
	var request createRolesRequest
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

	role := models.Roles{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		PhoneNumber:  request.PhoneNumber,
		Birthday:     request.Birthday,
	}

	id, err := h.rolesRepo.Create(c, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not create role"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Update godoc
// @Tags         roles - это выполняет роль Админа
// @Summary      Update role
// @Accept       json
// @Produce      json
// @Param id path int true "Role id"
// @Param request body handlers.updateRolesRequest true "Role data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "Role not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /roles/{id} [put]
func (h *RolesHandlers) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid role Id"))
		return
	}

	var request updateRolesRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	roles, err := h.rolesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("Role not found"))
		return
	}

	roles.Name = request.Name
	roles.Email = request.Email
	roles.PhoneNumber = request.PhoneNumber
	roles.Birthday = request.Birthday

	err = h.rolesRepo.Update(c, id, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Tags         roles - это выполняет роль Админа
// @Summary      Delete role
// @Accept       json
// @Produce      json
// @Param id path int true "Role id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "Role not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /roles/{id} [delete]
func (h *RolesHandlers) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid role Id"))
		return
	}

	_, err = h.rolesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("Role not found"))
		return
	}

	err = h.rolesRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// ChangePassword godoc
// @Tags         roles - это выполняет роль Админа
// @Summary      Change role password
// @Accept       json
// @Produce      json
// @Param id path int true "Role id"
// @Param request body handlers.changePasswordRequest true "Password data"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "Role not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /roles/{id}/changePassword [patch]
func (h *RolesHandlers) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid role Id"))
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

	role, err := h.rolesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("Role not found"))
		return
	}

	role.PasswordHash = string(passwordHash)

	err = h.rolesRepo.Update(c, id, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// FindById godoc
// @Tags         поиск user по id(Админ ищет юзера)
// @Summary      Role Finds user by id
// @Accept       json
// @Produce      json
// @Param id path int true "Roles id"
// @Success      200  {array}  handlers.rolesResponse  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid role id"
// @Failure   	 404  {object} models.ApiError "Role not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /rolesusers/{id} [get]
func (h *RolesHandlers) FindByIdUser(c *gin.Context) {
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
// @Tags 		 Админ получает список юзеров
// @Summary      Roles gets user's list
// @Produce      json
// @Success      200  {array} handlers.userResponse "OK"
// @Failure   	 500  {object} models.ApiError
// @Router       /rolesuser [get]
func (h *RolesHandlers) FindAllUsers(c *gin.Context) {
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

// Update godoc
// @Tags         Админ редактирует данные юзера
// @Summary      Role Updates user´s information
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Param request body handlers.updateUserRequest true "User data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /rolesuser/{id} [put]
func (h *RolesHandlers) UpdateUser(c *gin.Context) {
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
// @Tags         Админ удаляет юзера
// @Summary      Delete user
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router      /rolesuser/{id} [delete]
func (h *RolesHandlers) DeleteUser(c *gin.Context) {
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
// @Tags         Админ меняет пароль юзера
// @Summary      Change user password
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Param request body handlers.changePasswordRequest true "Password data"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /rolesuser/{id}/changePassword [patch]
func (h *RolesHandlers) ChangePasswordUser(c *gin.Context) {
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

// FindById godoc
// @Summary      Find by id
// @Tags         Админ ищет фильм
// @Accept       json
// @Produce      json
// @Param        id path int true "Movie id"
// @Success      200  {object}  models.MovieAdminResponse
// @Failure      400  {object}  models.ApiError "Invalid Movie Id"
// @Router       /rolesmovie/{id} [get]
func (h *RolesHandlers) FindByIdMoviesAdmin(c *gin.Context) {
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
// @Tags         Админ получает список фильмов
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.MovieAdminResponse "List of movies"
// @Failure      500  {object}  models.ApiError "Internal Server Error"
// @Router       /rolesmovie [get]
func (h *RolesHandlers) FindAllMoviesforAdmin(c *gin.Context) {
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

func (h *RolesHandlers) saveMoviesPoster(c *gin.Context, poster *multipart.FileHeader) (string, error) {
	filename := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(poster.Filename))
	filepath := fmt.Sprintf("images/%s", filename)
	err := c.SaveUploadedFile(poster, filepath)

	return filename, err
}

// Update godoc
// @Summary      Update moviesAdmin
// @Tags         Админ редактирует фильм
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
// @Router       /rolesmovie/{id} [put]
func (h *RolesHandlers) UpdateMovies(c *gin.Context) {
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
// @Tags        Админ удаляет фильм
// @Accept       json
// @Produce      json
// @Param        id path int true "Movie id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /rolesmovie/{id} [delete]
func (h *RolesHandlers) DeleteMovie(c *gin.Context) {
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
// @Tags         количество просмотров доступен только для Админа
// @Accept       json
// @Produce      json
// @Param id path int true "Movie id"
// @Param isWatched query bool true "Flag value"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /rolesmovie/{id}/setWatched [patch]
func (h *RolesHandlers) HandleSetWatched(c *gin.Context) {
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
