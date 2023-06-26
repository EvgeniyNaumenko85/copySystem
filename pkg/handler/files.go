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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": err.Error(),
		})
		return
	}

	filePath, err := h.services.GetFileByID(fileID, userName)
	if err != nil {
		switch err {
		case models.ErrFileAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{
				"reason": err.Error(),
			})
			return
		case models.ErrFileNotExists:
			c.JSON(http.StatusNoContent, gin.H{})
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
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))

	c.File(filePath)

	return
}

func (h *Handler) allFilesInfo(c *gin.Context) {
	files, err := h.services.AllFilesInfo()

	if err != nil {
		switch err {
		case models.ErrFilesNotExists:
			c.JSON(http.StatusNoContent, gin.H{})
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

func (h *Handler) showAllUserFilesInfo(c *gin.Context) {

	files, err := h.services.ShowAllUserFilesInfo(c)

	if err != nil {
		switch err {
		case models.ErrNoRows:
			c.JSON(http.StatusNoContent, gin.H{})
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
			c.JSON(http.StatusNoContent, gin.H{})
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
		switch err {
		case models.ErrNoRows:
			c.JSON(http.StatusNoContent, gin.H{})
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

	c.JSON(http.StatusOK, "file successfully deleted")
}

func (h *Handler) deleteAllFiles(c *gin.Context) {
	err := h.services.DeleteAllFiles()
	if err != nil {
		switch err {
		case models.ErrNoRows:
			c.JSON(http.StatusNoContent, gin.H{})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while deleting user",
				"reason":  err,
			})
		}
		return
	}

	c.JSON(http.StatusOK, "all files successfully deleted")
}
