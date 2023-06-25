package handler

import (
	"copySys/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userRoleCtx         = "userRole"
	userNameCtx         = "userName"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "empty auth header"})

		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "invalid auth header"})

		return
	}

	if len(headerParts[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "token is empty"})

		return
	}

	userId, role, userName, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": err.Error()})

		return
	}

	c.Set(userCtx, userId)
	c.Set(userRoleCtx, role)
	c.Set(userNameCtx, userName)
}

func getUserRole(c *gin.Context) (string, error) {
	roleCtx, ok := c.Get(userRoleCtx)

	if !ok {
		//todo создать новую ошибку в models
		logger.Error.Println("user role not found")
		return "", errors.New("user role not found")
	}

	role, ok := roleCtx.(string)
	if !ok {
		//todo создать новую ошибку в models
		logger.Error.Println("user role is of invalid type")
		return "", errors.New("user role is of invalid type")
	}

	return role, nil
}

func IdentifyUserRole(c *gin.Context) {
	role, _ := getUserRole(c)

	if role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"reason": "access denied"})
		return
	}
}

func IdMiddleware(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error.Println("Hello from middleware: ", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "invalid id",
		})
		return
	}

	c.Next()
}
