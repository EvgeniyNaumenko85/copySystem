package utils

import (
	"copySys/models"
	"copySys/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
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
