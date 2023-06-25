package handler

import (
	"copySys/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("/api")

	files := api.Group("/files", h.userIdentity)
	{
		files.POST("/", h.uploadFile)
		files.GET("/:id", IdMiddleware, h.getFileByID)
		files.GET("/", h.showAllUserFilesInfo)
		files.GET("/all", IdentifyUserRole, h.allFilesInfo)
		files.POST("/name", h.findFileByFileName)
		files.DELETE("/:id", IdMiddleware, h.deleteFileByID)
		files.DELETE("/all", IdentifyUserRole, h.deleteAllFiles)
	}

	access := api.Group("/access", h.userIdentity)
	{
		access.POST("/", h.providingAccess)
		access.POST("/:id", IdMiddleware, h.providingAccessAll)
		access.DELETE("/", h.removeAccess)
		//todo роут на ограничение прав всем пользоватевалям кроме владельца
		//access.DELETE("/all", h.stopAccess) // отправляем JSON c парой file_id/user_id для удаления из таблицы access
	}

	//todo //stat := api.Group("/", h.userIdentity)
	// роут групп по получению статистики по файлам (только админу?)
	{
		//todo роут на получение информации о типе, кол-ве, общем объеме файлов конкретного пользователя (возможно
		// добавить столбец в таблицу files с инфой о размере каждого файла, которая вносится в нее при записи файла в папку)
		//stat.GET("/:id", IdMiddleware, IdentifyUserRole, h.getStat)
		//todo роут на получение статистики всех пользователей (только для админа)
		//stat.GET("/all", IdMiddleware, IdentifyUserRole, h.allStat)
	}

	users := api.Group("/users", h.userIdentity)
	{
		users.GET("", h.getAllUsers)
		users.GET("/:id", IdMiddleware, h.getUserByID)
		users.PUT("/:id", IdentifyUserRole, IdMiddleware, h.updateUserByID)
		users.DELETE("/:id", IdentifyUserRole, IdMiddleware, h.deleteUserByID)
	}

	auth := api.Group("/auth")
	{
		auth.POST("signUp", h.signUp)
		auth.POST("signIn", h.signIn)
	}

	router.GET("/", Ping)

	return router
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"reason": "up and working",
	})
}
