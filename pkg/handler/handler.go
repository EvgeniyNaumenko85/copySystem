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
	router.GET("/", Ping)

	api := router.Group("/api")
	//api := router.Group("/auth")

	tasks := api.Group("/tasks", h.userIdentity)
	{
		tasks.GET("/", h.getAllTasks)
		tasks.GET("/:id", IdMiddleware, h.getTask)
		tasks.POST("/", h.createTask)
		tasks.PUT("/:id", IdMiddleware, h.updateTask)
		tasks.PUT("/reassign/:id", IdMiddleware, h.reassignTask)
		tasks.DELETE("/:id", IdentifyUserRole, IdMiddleware, h.deleteTask)
		tasks.GET("/:id/overdue", IdMiddleware, h.getOverdueTasks)
		tasks.GET("/:id/tasks", IdMiddleware, h.getTaskByUserId)
		tasks.GET("/:id/undone_tasks", IdMiddleware, h.getUndoneTasksByUserId)
	}

	users := api.Group("/users", h.userIdentity)
	{
		users.GET("", h.getAllUsers)
		users.GET("/:id", IdentifyUserRole, IdMiddleware, h.getUser)
		users.PUT("/:id", IdentifyUserRole, IdMiddleware, h.updateUser)
		users.DELETE("/:id", IdentifyUserRole, IdMiddleware, h.deleteUser)
	}

	auth := api.Group("/auth")
	{
		auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
	}

	return router
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"reason": "up and working",
	})
}
