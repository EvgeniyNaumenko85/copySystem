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

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func getUserRole(c *gin.Context) (string, error) {
	roleCtx, ok := c.Get(userRoleCtx)

	if !ok {
		logger.Error.Println("user role not found")
		return "", errors.New("user role not found")
	}

	role, ok := roleCtx.(string)
	if !ok {
		logger.Error.Println("user role is of invalid type")
		return "", errors.New("user role is of invalid type")
	}

	return role, nil
}

func IdentifyUserRole(c *gin.Context) {
	//id, _ := getUserId(c)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid id",
		})
		return
	}

	c.Next()
}

/*
func FileSizeMiddleware(c *gin.Context) {
	fmt.Println("Hello from SizeMiddleware")

	// to do: внедрить регулирование ограничения объема файла  (передача параметра по роуту)
	// Ограничение размера файла до 10 МБ (10 * 1024 * 1024 байт)
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Error parsing multipart form:", err)
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"err": "Error parsing multipart form",
		})
		return
	}

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file:", err)
		c.String(http.StatusBadRequest, "Error retrieving file")
		return
	}
	defer file.Close()

	// Проверка размера файла
	if handler.Size > 10<<20 {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"err": "File size exceeds the limit",
		})

	}

	c.Next()
}
*/
