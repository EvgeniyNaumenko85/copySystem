package handler

import (
	"copySys/models"
	"copySys/pkg/logger"
	"copySys/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

func (h *Handler) uploadFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "error getting file"})
		return
	}

	fileId, err := h.services.UploadFile(header, c)
	if err != nil {
		switch err {
		case models.ErrFileToBig:
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"reason": "file to upload is too big",
				"err":    err.Error(),
			})
			return
		case models.ErrFileAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"reason": "file already exists",
				"err":    err.Error(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"reason": "error while saving file to db",
				"err":    err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file saved successfully",
		"id":      fileId,
	})
}

/*
func (h *Handler) getFileByID(c *gin.Context) {
	idStr := c.Param("id")
	fileId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid id",
		})
		return
	}

	err = h.services.GetFileByID(fileId, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "error while load file from db",
			"err":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file loaded successfully"})
}
*/

func (h *Handler) getFileByID(c *gin.Context) {
	idStr := c.Param("id")
	fileID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid id",
		})
		return
	}

	userName, err := utils.GetUserNameFromContext(c)

	filePath, err := h.services.GetFileByID(fileID, userName)
	if err != nil {
		switch err {
		case models.ErrFileAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{
				"reason": err.Error(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"reason": err.Error(),
			})
			return
		}
	}

	// Устанавливаем заголовки для скачивания файла
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	// Устанавливаем заголовок Content-Disposition для передачи имени файла с расширением
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))

	c.File(filePath)

	return
}

func (h *Handler) allFilesInfo(c *gin.Context) {
	files, err := h.services.AllFilesInfo()

	if err != nil {
		switch err {
		case models.ErrNoRows:
			c.JSON(http.StatusNoContent, gin.H{
				"reason": err.Error(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"reason": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, files)
}

// todo проверить работу:
func (h *Handler) showAllUserFilesInfo(c *gin.Context) {

	files, err := h.services.ShowAllUserFilesInfo(c)

	if err != nil {
		switch err {
		case models.ErrNoRows:
			c.JSON(http.StatusNoContent, gin.H{
				"reason": err.Error(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"reason": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, files)
}

// todo
func (h *Handler) findFileByFileName(c *gin.Context) {
	var f models.File
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "error while binding body",
		})
		return
	}

	userName, err := utils.GetUserNameFromContext(c)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	file, err := h.services.FindFileByFileName(f.FileName, userName)
	if err != nil {
		switch err {
		case models.ErrNoRows:
			c.JSON(http.StatusNoContent, gin.H{
				"reason": err,
			})
			return
		case models.ErrFileAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{
				"reason": err.Error(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"reason": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, file)
}

func (h *Handler) deleteFileByID(c *gin.Context) {
	idStr := c.Param("id")
	fileID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "invalid file id",
		})
		return
	}

	err = h.services.DeleteFileByID(fileID)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "file is not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while deleting file",
				"reason":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, "file successfully deleted")
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
