package handler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	if !h.isAllowedFileType(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file type is not allowed",
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

	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to prepare temporary file",
		})
		return
	}

	defer os.Remove(tempPath)

	if err := c.SaveUploadedFile(
		file,
		tempPath,
	); err != nil {
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
		errorResponse(
			c,
			http.StatusBadRequest,
			"file is required",
		)
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
			errorResponse(
				c,
				http.StatusBadRequest,
				"file is required",
			)
			return
		}

		if errors.Is(err, service.ErrFileExpired) {
			errorResponse(
				c,
				http.StatusBadRequest,
				"file is required",
			)
			return
		}

		errorResponse(
			c,
			http.StatusBadRequest,
			"file is required",
		)

		return
	}

	c.Header(
		"Content-Disposition",
		`attachment; filename="`+file.OriginalName+`"`,
	)

	c.File(file.FilePath)
}

func (h *FileHandler) isAllowedFileType(
	filename string,
) bool {

	extension := strings.TrimPrefix(
		strings.ToLower(filepath.Ext(filename)),
		".",
	)

	for _, allowed := range h.allowedFileTypes {

		if extension == strings.ToLower(
			strings.TrimSpace(allowed),
		) {
			return true
		}
	}

	return false
}
