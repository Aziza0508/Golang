package handlers

import (
	"encoding/json"
	"net/http"
	"practice2/internal/models"
	"practice2/internal/storage"
	"strconv"
	"strings"
)

type TaskHandler struct {
	storage *storage.TaskStorage
}

func NewTaskHandler(s *storage.TaskStorage) *TaskHandler {
	return &TaskHandler{storage: s}
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "method not allowed"})
	}
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		tasks := h.storage.GetAll()
		json.NewEncoder(w).Encode(tasks)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid id"})
		return
	}

	task, exists := h.storage.GetByID(id)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "task not found"})
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid title"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid title"})
		return
	}

	task := h.storage.Create(title)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid id"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid id"})
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Done == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid done"})
		return
	}

	_, exists := h.storage.Update(id, *req.Done)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "task not found"})
		return
	}

	json.NewEncoder(w).Encode(models.UpdateResponse{Updated: true})
}
