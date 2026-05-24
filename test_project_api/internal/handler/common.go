package handler

import (
    "encoding/json"
    "net/http"
    apperrors "org-api/internal/errors"
    log "github.com/sirupsen/logrus"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Error("Failed to encode JSON: ", err)
    }
}

func writeError(w http.ResponseWriter, err error) {
    if appErr, ok := err.(*apperrors.AppError); ok {
        http.Error(w, appErr.Message, appErr.Code)
        return
    }
    log.Error("Unexpected error: ", err)
    http.Error(w, "Internal server error", http.StatusInternalServerError)
}