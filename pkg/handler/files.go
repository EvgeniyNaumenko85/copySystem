package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) uploadFile(c *gin.Context) {
	//file, err := c.FormFile("file")
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при получении файла"})
		return
	}

	err = h.services.UploadFile(file, header, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "error while saving file to db",
			"err":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно сохранен"})
}

func (h *Handler) getFile(c *gin.Context) {
	fmt.Println("Hello from loadFile")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid id",
		})
		return
	}

	err = h.services.GetFile(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "error while load file from db",
			"err":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно выгружен"})
}

/*
func (h *Handler) createTask(c *gin.Context) {
	var t *models.Task
	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
		})
		return
	}

	id, err := h.services.CreateTask(*t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "error while saving task to db",
			//"err":    err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"reason": "successfully created",
		"id":     id,
	})
}

func (h *Handler) getAllTasks(c *gin.Context) {
	t, err := h.services.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, t)

}

func (h *Handler) getTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid task id",
		})
		return
	}

	tasks, err := h.services.GetTaskByID(id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "task is not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while getting task",
				"reason":  err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) getOverdueTasks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	t, err := h.services.GetOverdueTasks(id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "task(s) is not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while getting task",
				"reason":  err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) updateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid task id",
		})
		return
	}

	var t *models.Task
	if err = c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
		})
		return
	}

	if err = h.services.UpdateTaskByID(id, *t); err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while updating task",
				"reason":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reason": "successfully updated",
	})
}

func (h *Handler) reassignTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	var r *models.UserRequest
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
		})
		return
	}

	oldUserID := r.OldUserId
	newUserID := r.NewUserId

	if err = h.services.ReassignTask(oldUserID, newUserID, id); err != nil {
		switch err.Error() {
		case fmt.Sprintf("pq: Задача %d не найдена для пользователя %d", id, oldUserID):
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
				"reason":  err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to reassign task",
				"reason":  err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Task reassigned successfully",
	})
}

func (h *Handler) deleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid task id",
		})
		return
	}

	err = h.services.DeleteTaskByID(id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "task is not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while updating task",
				"reason":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, "task successfully deleted")
}

func (h *Handler) getTaskByUserId(c *gin.Context) {
	userId, err := getUserId(c)
	//idStr := c.Param("id")
	//id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid task id",
		})
		return
	}

	task, err := h.services.GetTaskByUserID(userId)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while getting task",
				"reason":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) getUndoneTasksByUserId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid task id",
		})
		return
	}

	task, err := h.services.GetUndoneTasksByUserID(id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "tasks is not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while getting task",
				"reason":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}
*/