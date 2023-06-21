package handler

import (
	"copySys/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var payload models.User

	err := c.ShouldBind(&payload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	// Hashing password
	userId, err := h.services.CreateUser(payload)
	if err != nil {
		if err.Error() == models.ErrNotUnicUser {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": userId})
}

type singInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input singInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	// Hashing password
	token, err := h.services.GenerateToken(input.UserName, input.Password, input.Role)
	if err != nil {

		if err.Error() == models.ErrNoRowsSQL {
			c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
