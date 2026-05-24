package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "org-api/internal/service"
    apperrors "org-api/internal/errors"
)

type DepartmentHandler struct {
    service *service.DepartmentService
}

func NewDepartmentHandler(s *service.DepartmentService) *DepartmentHandler {
    return &DepartmentHandler{service: s}
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name     string `json:"name"`
        ParentID *uint  `json:"parent_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, apperrors.ErrInvalidInput)
        return
    }
    dept, err := h.service.Create(req.Name, req.ParentID)
    if err != nil {
        writeError(w, err)
        return
    }
    writeJSON(w, http.StatusCreated, dept)
}

func (h *DepartmentHandler) Get(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        writeError(w, apperrors.ErrInvalidInput)
        return
    }
    depthStr := r.URL.Query().Get("depth")
    depth := 1
    if depthStr != "" {
        d, err := strconv.Atoi(depthStr)
        if err == nil && d >= 0 {
            depth = d
        }
    }
    includeEmployees := true
    if inc := r.URL.Query().Get("include_employees"); inc == "false" {
        includeEmployees = false
    }
    dept, err := h.service.GetByID(uint(id), depth, includeEmployees)
    if err != nil {
        writeError(w, err)
        return
    }
    writeJSON(w, http.StatusOK, dept)
}

func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        writeError(w, apperrors.ErrInvalidInput)
        return
    }
    var req struct {
        Name     *string `json:"name"`
        ParentID *uint   `json:"parent_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, apperrors.ErrInvalidInput)
        return
    }
    dept, err := h.service.Update(uint(id), req.Name, req.ParentID)
    if err != nil {
        writeError(w, err)
        return
    }
    writeJSON(w, http.StatusOK, dept)
}

func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        writeError(w, apperrors.ErrInvalidInput)
        return
    }
    mode := r.URL.Query().Get("mode")
    var reassignTo *uint
    if mode == "reassign" {
        val := r.URL.Query().Get("reassign_to_department_id")
        if val == "" {
            writeError(w, apperrors.ErrInvalidInput)
            return
        }
        i, err := strconv.ParseUint(val, 10, 32)
        if err != nil {
            writeError(w, apperrors.ErrInvalidInput)
            return
        }
        reassignTo = new(uint)
        *reassignTo = uint(i)
    } else if mode != "cascade" {
        writeError(w, apperrors.ErrInvalidInput)
        return
    }
    if err := h.service.Delete(uint(id), mode, reassignTo); err != nil {
        writeError(w, err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}