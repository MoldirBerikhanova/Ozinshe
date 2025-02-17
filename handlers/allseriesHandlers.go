package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

type AllSeriesHandlers struct {
	allseriesRepo *repositories.AllSeriesRepository
}

type createAllSeriesRequest struct {
	Series      *int    `form:"series"`
	Title       *string `form:"title"`
	Description *string `form:"description"`
	ReleaseYear *int    `form:"release_year"`
	Director    *string `form:"director"`
	Rating      *int    `form:"rating"`
	TrailerUrl  *string `form:"trailer_url"`
}

type updateAllSeriesRequest struct {
	Series      *int    `form:"series"`
	Title       *string `form:"title"`
	Description *string `form:"description"`
	ReleaseYear *int    `form:"release_year"`
	Director    *string `form:"director"`
	Rating      *int    `form:"rating"`
	TrailerUrl  *string `form:"trailer_url"`
}

func NewAllSeriesHandlers(allseriesRepo *repositories.AllSeriesRepository) *AllSeriesHandlers {
	return &AllSeriesHandlers{
		allseriesRepo: allseriesRepo,
	}
}

// Create godoc
// @Summary      Create allseries
// @Tags 		 allseries - это эндпоинты для каждой серии
// @Accept       json
// @Produce      json
// @Param request body models.AllSeries true "AllSeries model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid request AllSeries"
// @Failure   	 500  {object} models.ApiError
// @Router       /allseries [post]
func (h *AllSeriesHandlers) Create(c *gin.Context) {
	var request createAllSeriesRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request AllSeries")
		return
	}

	allserie := models.AllSeries{
		Series:      request.Series,
		Title:       request.Title,
		Description: request.Description,
		ReleaseYear: request.ReleaseYear,
		Director:    request.Director,
		Rating:      request.Rating,
		TrailerUrl:  request.TrailerUrl,
	}

	id, err := h.allseriesRepo.Create(c, allserie)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// FindById godoc
// @Summary      Find by id allseries
// @Tags         allseries - это эндпоинты для каждой серии
// @Accept       json
// @Produce      json
// @Param        id path int true "AllSeries id"
// @Success      200  {object}  models.AllSeries "Ok"
// @Failure      400  {object}  models.ApiError "Invalid allseries id"
// @Router       /allseries/{id} [get]
func (h *AllSeriesHandlers) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid allseries id"))
		return
	}

	_, err = h.allseriesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

}

// FindAll godoc
// @Summary      Get all allseries
// @Tags          allseries - это эндпоинты для каждой серии
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.AllSeries "List of allseries"
// @Failure      500  {object}  models.ApiError "Internal Server Error"
// @Router       /allseries [get]
func (h *AllSeriesHandlers) FindAll(c *gin.Context) {
	allseries, err := h.allseriesRepo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, allseries)
}

// Update godoc
// @Summary      Update allseries
// @Tags 		 allseries - это эндпоинты для каждой серии
// @Accept       json
// @Produce      json
// @Param request body models.AllSeries true "AllSeries model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid AllSeries Id"
// @Failure   	 500  {object} models.ApiError
// @Router       /allseries/{id} [put]
func (h *AllSeriesHandlers) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid AllSeries Id"))
		return
	}

	_, err = h.allseriesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request updateAllSeriesRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	allserie := models.AllSeries{
		Series:      request.Series,
		Title:       request.Title,
		Description: request.Description,
		ReleaseYear: request.ReleaseYear,
		Director:    request.Director,
		Rating:      request.Rating,
		TrailerUrl:  request.TrailerUrl,
	}

	err = h.allseriesRepo.Update(c, id, allserie)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary      Delete allseries
// @Tags          allseries - это эндпоинты для каждой серии
// @Accept       json
// @Produce      json
// @Param        id path int true "Allseries id"
// @Success      200  {object}  models.AllSeries "Ok"
// @Failure      400  {object}  models.ApiError "Invalid AllSeries Id"
// @Router       /allseries/{id} [delete]
func (h *AllSeriesHandlers) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid AllSeries Id"))
		return
	}

	_, err = h.allseriesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = h.allseriesRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
