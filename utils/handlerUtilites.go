package utils

import (
	"copySys/models"
	"copySys/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
)

func GetUserNameFromContext(c *gin.Context) (string, error) {
	userNameTypeAny, ok := c.Get("userName")
	if !ok {
		logger.Error.Println(models.ErrCantGetUserName.Error())
		return "", models.ErrCantGetUserName
	} else {
		userName := fmt.Sprintf("%v", userNameTypeAny)
		return userName, nil
	}
}

func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}
