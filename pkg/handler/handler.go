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
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("signUp", h.signUp)
		auth.POST("signIn", h.signIn)
	}

	users := api.Group("/users", h.userIdentity)
	{
		users.GET("", h.getAllUsers)
		users.GET("/:id", IdMiddleware, h.getUserByID)
		users.PUT("/:id", IdentifyUserRole, IdMiddleware, h.updateUserByID)
		users.DELETE("/:id", IdentifyUserRole, IdMiddleware, h.deleteUserByID)
	}

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
		access.DELETE("/all", h.removeAccessToAll)
	}

	limits := api.Group("/limits", h.userIdentity)
	{
		limits.PUT("/:id", IdentifyUserRole, IdMiddleware, h.setLimitsToUser)
	}

	stat := api.Group("/stat", h.userIdentity)
	{
		stat.GET("/:id", IdMiddleware, IdentifyUserRole, h.getUserStatistics)
	}

	router.GET("/", Ping)

	return router
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"reason": "up and working",
	})
}
