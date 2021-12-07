package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/FlameInTheDark/file-share/app/api/handlers/requests"
	"github.com/FlameInTheDark/file-share/app/api/handlers/responses"
	"github.com/FlameInTheDark/file-share/business/service/files"
	"github.com/FlameInTheDark/file-share/business/service/storage"
	"github.com/FlameInTheDark/file-share/foundation/s3"
)

type Handler struct {
	file    FileService
	storage StorageService
}

func NewHandler(db *sqlx.DB, s3storage *s3.MinioClient) *Handler {
	handler := Handler{
		file:    files.NewFilesService(db),
		storage: storage.NewStorageService(s3storage),
	}
	return &handler
}

func (h *Handler) Download(c echo.Context) error {
	id, name, err := h.file.Find(c.Request().Context(), c.Param("id"))
	if err != nil {
		return responses.Error(c, err)
	}

	url, err := h.storage.Download(c.Request().Context(), id, name)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(http.StatusOK, responses.DownloadResponse{URL: url})
}

func (h *Handler) Upload(c echo.Context) error {
	var req requests.UploadRequest
	err := c.Bind(&req)
	if err != nil {
		return responses.Error(c, err)
	}

	id, err := h.file.Create(c.Request().Context(), req.Name)
	if err != nil {
		return responses.Error(c, err)
	}

	url, err := h.storage.Upload(c.Request().Context(), id)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(http.StatusOK, responses.UploadResponse{URL: url, ID: id})
}

func (h *Handler) Statistics(c echo.Context) error {
	stats, err := h.file.Statistics(c.Request().Context(), c.Param("id"))
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(http.StatusOK, responses.StatisticsResponse{Downloads: stats})
}