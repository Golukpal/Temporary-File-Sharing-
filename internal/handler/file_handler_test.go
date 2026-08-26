package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Golukpal/file-sharing/internal/service"
)

func TestDownload_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	fileService := service.NewFileService(nil)

	fileHandler := NewFileHandler(
		fileService,
		10*1024*1024,
		[]string{"pdf", "txt"},
	)

	router.GET(
		"/files/:id",
		fileHandler.Download,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/files/hello",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		rec.Code,
	)
}
