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

/*
func GetInfoFromContext(headerKey string, c *gin.Context) (string, error) {
	fmt.Println("headerKey: ", headerKey)

	headerKey = "userName"

	fileNameAny := c.GetHeader(headerKey)
	fmt.Println("fileNameAny: ", fileNameAny)

	//if !ok {
	//	fmt.Println("fileNameAny: ", fileNameAny)
	//	logger.Error.Println(models.ErrCantGetInfoFromHeader.Error())
	//	return "", models.ErrCantGetInfoFromHeader
	//} else {
	//	fileName, success := fileNameAny.(string)
	//	if !success {
	//		return "", errors.New("value is not a string")
	//	}
	//	return fileName, nil
	//}

	return fileNameAny, nil
}
*/
