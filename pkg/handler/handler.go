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

	files := api.Group("/files", h.userIdentity)

	{
		//files.POST("/", FileSizeMiddleware, h.uploadFile)
		//todo не забывать редактироваать соответсвующие таблицы!!!!
		files.POST("/", h.uploadFile)
		files.GET("/:id", IdMiddleware, h.getFileByID)
		files.GET("/all", IdentifyUserRole, h.allFilesInfo)
		files.GET("/", h.showAllUserFilesInfo)
		files.DELETE("/:id", IdMiddleware, h.deleteFileByID)

		//todo роут под поиска файла пользователем
		// files.GET("/file/", h.findFile)
		//(отправляется JSON c именем файла, распарсивается имя требуемого файла, затем ищем файл в таблице files,
		//возвращаем id файла из таблицы files)

		//роуты под рефакторинг
		//files.GET("/", h.getAllTasks)
		//files.GET("/:id", IdMiddleware, h.getTask)
		//files.POST("/", h.createTask)
		////files.PUT("/:id", IdMiddleware, h.updateTask)
		//files.PUT("/reassign/:id", IdMiddleware, h.reassignTask)
		//files.DELETE("/:id", IdentifyUserRole, IdMiddleware, h.deleteTask)
		//files.GET("/:id/overdue", IdMiddleware, h.getOverdueTasks)
		//files.GET("/:id/tasks", IdMiddleware, h.getTaskByUserId)
		//files.GET("/:id/undone_tasks", IdMiddleware, h.getUndoneTasksByUserId)
	}
	/*
		//todo //access := api.Group("/", h.userIdentity)
		//  роут групп по управлению доступом к файлам
		{
			//todo роут на добавление прав другому пользоватевалю
			//access.POST("/", h.providingAccess) // отправляем JSON c парой file_id/user_id для добавления в таблицу access
			//todo роут на добавление прав всем пользоватевалям
			//access.POST("/all", h.fullAccess) // выполняется функция, организующая запись всех возможных пар
			//file_id/user_id для добавления в таблицу access
			//todo роут на удаление прав пользоватеваля
			//access.DELETE("/", h.limitationAccess) // отправляем JSON c парой file_id/user_id для удаления из таблицы access
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

	*/

	users := api.Group("/users", h.userIdentity)
	{
		users.GET("", h.getAllUsers)
		users.GET("/:id", IdentifyUserRole, IdMiddleware, h.getUserByID)
		users.PUT("/:id", IdentifyUserRole, IdMiddleware, h.updateUserByID)
		users.DELETE("/:id", IdentifyUserRole, IdMiddleware, h.deleteUserByID)
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
