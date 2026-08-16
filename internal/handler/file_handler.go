package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Golukpal/file-sharing/internal/service"
)

type FileHandler struct {
	service *service.FileService
}

func NewFileHandler(service *service.FileService) *FileHandler {
	return &FileHandler{
		service: service,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})
		return
	}

	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file cannot be empty",
		})
		return
	}

	tempPath := file.TempFile

	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save uploaded file",
		})
		return
	}

	result, err := h.service.Upload(
		c.Request.Context(),
		file.Filename,
		file.Size,
		tempPath,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}