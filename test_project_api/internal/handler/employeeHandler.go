package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "github.com/gorilla/mux"
    "org-api/internal/service"
    apperrors "org-api/internal/errors"
)

type EmployeeHandler struct {
    service *service.EmployeeService
}

func NewEmployeeHandler(s *service.EmployeeService) *EmployeeHandler {
    return &EmployeeHandler{service: s}
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    deptID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        writeError(w, apperrors.ErrInvalidInput)
        return
    }
    var req struct {
        FullName string     `json:"full_name"`
        Position string     `json:"position"`
        HiredAt  *time.Time `json:"hired_at"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, apperrors.ErrInvalidInput)
        return
    }
    emp, err := h.service.Create(uint(deptID), req.FullName, req.Position, req.HiredAt)
    if err != nil {
        writeError(w, err)
        return
    }
    writeJSON(w, http.StatusCreated, emp)
}