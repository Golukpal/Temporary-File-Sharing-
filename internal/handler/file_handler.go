package handler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/Golukpal/file-sharing/internal/service"
)

type FileHandler struct {
	service          *service.FileService
	maxFileSize      int64
	allowedFileTypes []string
}

func NewFileHandler(
	service *service.FileService,
	maxFileSize int64,
	allowedFileTypes []string,
) *FileHandler {
	return &FileHandler{
		service:          service,
		maxFileSize:      maxFileSize,
		allowedFileTypes: allowedFileTypes,
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

	if file.Size > h.maxFileSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "file size exceeds maximum allowed size",
		})
		return
	}

	tempFile, err := os.CreateTemp(
		"",
		"upload-*"+filepath.Ext(file.Filename),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create temporary file",
		})
		return
	}

	tempPath := tempFile.Name()
	tempFile.Close()

	defer os.Remove(tempPath)

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

func (h *FileHandler) Download(c *gin.Context) {
	id := c.Param("id")

	file, err := h.service.GetFile(
		c.Request.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file not found",
			})
			return
		}

		if errors.Is(err, service.ErrFileExpired) {
			c.JSON(http.StatusGone, gin.H{
				"error": "file has expired",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get file",
		})
		return
	}

	c.Header(
		"Content-Disposition",
		`attachment; filename="`+file.OriginalName+`"`,
	)

	c.File(file.FilePath)
}