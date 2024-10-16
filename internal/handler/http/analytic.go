package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/himmel520/uoffer/require/internal/repository"
	"github.com/himmel520/uoffer/require/models"
)

// @Summary Добавить аналитику
// @Description Добавляет новую аналитику
// @Tags analytic
// @Accept json
// @Produce json
// @Param analytic body models.Analytic true "Данные аналитики"
// @Success 200 {object} models.Analytic
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 409 {object} errorResponse "Conflict"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /admin/analytic [post]
func (h *Handler) addAnalytic(c *gin.Context) {
	var analytic *models.Analytic
	if err := c.BindJSON(&analytic); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		return
	}

	newAnalytic, err := h.srv.AddAnalytic(c.Request.Context(), analytic)
	switch {
	case errors.Is(err, repository.ErrPostIDExist):
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrAnalyticDependencyNotFound):
		c.AbortWithStatusJSON(http.StatusConflict, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, newAnalytic)
}

// @Summary Обновить аналитику
// @Description Обновляет данные аналитики по ID
// @Tags analytic
// @Accept json
// @Produce json
// @Param id path int true "ID аналитики"
// @Param analytic body models.AnalyticUpdate true "Обновленные данные аналитики"
// @Success 200 {object} models.Analytic
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 404 {object} errorResponse "Not Found"
// @Failure 409 {object} errorResponse "Conflict"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /admin/analytic/{id} [put]
func (h *Handler) updateAnalytic(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var analytic *models.AnalyticUpdate
	if err := c.BindJSON(&analytic); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		return
	}

	if analytic.IsEmpty() {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"analytic has no changes"})
		return
	}

	newAnalytic, err := h.srv.UpdateAnalytic(c.Request.Context(), id, analytic)
	switch {
	case errors.Is(err, repository.ErrPostIDExist):
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrAnalyticNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrAnalyticDependencyNotFound):
		c.AbortWithStatusJSON(http.StatusConflict, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, newAnalytic)
}

// @Summary Удалить аналитику
// @Description Удаляет аналитику по ID
// @Tags analytic
// @Param id path int true "ID аналитики"
// @Success 204 "No Content"
// @Failure 404 {object} errorResponse "Not Found"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /admin/analytic/{id} [delete]
func (h *Handler) deleteAnalytic(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.srv.DeleteAnalytic(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrAnalyticNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
			return
		}
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Получить аналитику с словами по ID поста
// @Description Возвращает аналитику с словами по ID поста
// @Tags Analytic
// @Param id path int true "ID поста"
// @Success 200 {object} []models.AnalyticWithWords
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 404 {object} errorResponse "Not Found"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /analytic/post/{id} [get]
func (h *Handler) getAnalyticWithWordsByPostID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	role, ok := c.Keys["role"]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"not found role"})
		return
	}

	analytics, err := h.srv.GetAnalyticWithWords(c.Request.Context(), id, role.(string))
	if err != nil {
		if errors.Is(err, repository.ErrAnalyticNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
			return
		}
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
