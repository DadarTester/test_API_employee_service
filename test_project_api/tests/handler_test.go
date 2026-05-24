package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gorilla/mux"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "org-api/internal/models"
    "org-api/internal/repository"
    "org-api/internal/service"
    "org-api/internal/handler"
    "github.com/stretchr/testify/assert"
)

func setupTestDB() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&models.Department{}, &models.Employee{})
    return db
}

func TestCreateDepartment(t *testing.T) {
    db := setupTestDB()
    deptRepo := repository.NewDepartmentRepository(db)
    empRepo := repository.NewEmployeeRepository(db)
    deptService := service.NewDepartmentService(deptRepo, empRepo, db, 5)
    deptHandler := handler.NewDepartmentHandler(deptService)

    reqBody := `{"name":"IT", "parent_id": null}`
    req := httptest.NewRequest("POST", "/departments", bytes.NewBufferString(reqBody))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    router := mux.NewRouter()
    router.HandleFunc("/departments", deptHandler.Create).Methods("POST")
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusCreated, rr.Code)
    var dept models.Department
    err := json.Unmarshal(rr.Body.Bytes(), &dept)
    assert.NoError(t, err)
    assert.Equal(t, "IT", dept.Name)
}