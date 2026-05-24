package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    log "github.com/sirupsen/logrus"

    "org-api/internal/config"
    "org-api/internal/repository"
    "org-api/internal/service"
    "org-api/internal/handler"
    "org-api/internal/midware"
)

func main() {

    cfg := config.Load()
    log.SetFormatter(&log.JSONFormatter{})
    log.SetLevel(log.InfoLevel)

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
	
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }

    deptRepo := repository.NewDepartmentRepository(db)
    empRepo := repository.NewEmployeeRepository(db)

    deptService := service.NewDepartmentService(deptRepo, empRepo, db, cfg.MaxDepth)
    empService := service.NewEmployeeService(empRepo, deptRepo)

    deptHandler := handler.NewDepartmentHandler(deptService)
    empHandler := handler.NewEmployeeHandler(empService)

    r := mux.NewRouter()
    r.Use(middleware.Logging)
    r.HandleFunc("/departments", deptHandler.Create).Methods("POST")
    r.HandleFunc("/departments/{id:[0-9]+}", deptHandler.Get).Methods("GET")
    r.HandleFunc("/departments/{id:[0-9]+}", deptHandler.Update).Methods("PATCH")
    r.HandleFunc("/departments/{id:[0-9]+}", deptHandler.Delete).Methods("DELETE")
    r.HandleFunc("/departments/{id:[0-9]+}/employees", empHandler.Create).Methods("POST")

    port := ":8080"
    log.Infof("Server starting on %s", port)
    if err := http.ListenAndServe(port, r); err != nil {
        log.Fatal("Server failed: ", err)
    }
}