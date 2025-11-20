package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

var AllowedExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
	".svg":  true,
}

type StorageHandler struct {
	storageService port.StorageService
	app            *application.Application
}

func NewStorageHandler(storageService port.StorageService, app *application.Application) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
		app:            app,
	}
}

func (h *StorageHandler) UploadRiskTypeIcon(w http.ResponseWriter, r *http.Request) {
	h.uploadIcon(w, r, "risks/types", true)
}

func (h *StorageHandler) UploadRiskTopicIcon(w http.ResponseWriter, r *http.Request) {
	h.uploadIcon(w, r, "risks/topics", false)
}

func (h *StorageHandler) uploadIcon(w http.ResponseWriter, r *http.Request, folder string, isType bool) {
	entityID := r.URL.Query().Get("id")
	if entityID == "" {
		writeJSONError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	if err := r.ParseMultipartForm(MaxFileSize); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	if header.Size > MaxFileSize {
		writeJSONError(w, http.StatusBadRequest, "file size exceeds 10MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedExtensions[ext] {
		writeJSONError(w, http.StatusBadRequest, "invalid file type. Allowed: png, jpg, jpeg, webp, svg")
		return
	}

	filename := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType([]byte(header.Filename))
	}

	if err := h.storageService.Upload(r.Context(), filename, file, contentType); err != nil {
		slog.Error("failed to upload to storage", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}

	if isType {
		if err := h.app.RiskUseCase.UpdateRiskTypeIcon(r.Context(), entityID, filename); err != nil {
			slog.Error("failed to update risk type icon", "id", entityID, "error", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to update icon")
			return
		}
	} else {
		if err := h.app.RiskUseCase.UpdateRiskTopicIcon(r.Context(), entityID, filename); err != nil {
			slog.Error("failed to update risk topic icon", "id", entityID, "error", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to update icon")
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"path": filename,
		"url":  h.storageService.GetURL(filename),
	})
}

func (h *StorageHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		writeJSONError(w, http.StatusBadRequest, "path is required")
		return
	}

	path = strings.TrimPrefix(path, "/")

	reader, err := h.storageService.Download(r.Context(), path)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "file not found")
		return
	}
	defer reader.Close()

	ext := strings.ToLower(filepath.Ext(path))
	contentType := "application/octet-stream"
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".webp":
		contentType = "image/webp"
	case ".svg":
		contentType = "image/svg+xml"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000")

	_, err = io.Copy(w, reader)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to serve file")
		return
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
