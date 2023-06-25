package handler

import (
	"copySys/models"
	"copySys/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) providingAccess(c *gin.Context) {
	var a models.AccessRequest
	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    models.ErrCantGetUserID,
			"reason": err,
		})
		return
	}

	fileID := a.FileID
	accessToUserID := a.AccessToUserID

	err = h.services.ProvidingAccess(fileID, accessToUserID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "access added successfully")
}

func (h *Handler) providingAccessAll(c *gin.Context) {
	idStr := c.Param("id")
	fileID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid file id",
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    models.ErrCantGetUserID,
			"reason": err,
		})
		return
	}

	err = h.services.ProvidingAccessAll(userID, fileID)
	if err != nil {
		switch err {
		case models.ErrFileAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{
				"reason": err.Error(),
			})
			return
		case models.ErrAccessInfoNotFound:
			c.JSON(http.StatusNoContent, gin.H{})
			return
		case models.ErrFileNotExists:
			c.JSON(http.StatusNoContent, gin.H{})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"reason": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, "access added successfully")
}

func (h *Handler) removeAccess(c *gin.Context) {
	var a models.AccessRequest
	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    models.ErrCantGetUserID,
			"reason": err,
		})
		return
	}

	fileID := a.FileID
	accessToUserID := a.AccessToUserID

	err = h.services.RemoveAccess(fileID, accessToUserID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "access removed successfully")
}

func (h *Handler) removeAccessToAll(c *gin.Context) {
	var a models.AccessRequest
	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    models.ErrCantGetUserID,
			"reason": err,
		})
		return
	}

	fileID := a.FileID

	err = h.services.RemoveAccessToAll(fileID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "access to all removed successfully")
}
