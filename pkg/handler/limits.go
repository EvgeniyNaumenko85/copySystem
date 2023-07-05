package handler

import (
	"copySys/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) setLimitsToUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid user id",
		})
		return
	}

	var l models.LimitRequest
	if err = c.BindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
			"err":    err,
		})
		return
	}

	err = h.services.SetLimitsToUser(userID, l.FileSizeLim, l.StorageSizeLim)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "limit added successfully")
}
