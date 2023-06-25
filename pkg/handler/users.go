package handler

import (
	"copySys/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getAllUsers(c *gin.Context) {
	user, err := h.services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid task id",
		})
		return
	}

	user, err := h.services.GetUserByID(id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while getting task",
				"reason":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid user id",
		})
		return
	}

	var u *models.User
	if err = c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
		})
		return
	}

	if err = h.services.UpdateUserByID(id, *u); err != nil {
		if err == models.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"reason": "user not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"reason": "error while updating user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reason": "successfully updated",
	})
}

func (h *Handler) deleteUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid user id",
		})
		return
	}

	err = h.services.DeleteUserByID(id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while deleting user",
				"reason":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, "user successfully deleted")
}
