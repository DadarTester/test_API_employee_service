package service

import (
    "strings"
    "time"
    "org-api/internal/models"
    "org-api/internal/repository"
    apperrors "org-api/internal/errors"
)

type EmployeeService struct {
    empRepo  *repository.EmployeeRepository
    deptRepo *repository.DepartmentRepository
}

func NewEmployeeService(empRepo *repository.EmployeeRepository, deptRepo *repository.DepartmentRepository) *EmployeeService {
    return &EmployeeService{
        empRepo:  empRepo,
        deptRepo: deptRepo,
    }
}

func (s *EmployeeService) Create(departmentID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error) {
    fullName = strings.TrimSpace(fullName)
    position = strings.TrimSpace(position)
    if fullName == "" || len(fullName) > 200 || position == "" || len(position) > 200 {
        return nil, apperrors.ErrInvalidInput
    }
    // проверка существования подразделения
    if _, err := s.deptRepo.GetByID(departmentID); err != nil {
        return nil, err
    }
    emp := &models.Employee{
        DepartmentID: departmentID,
        FullName:     fullName,
        Position:     position,
        HiredAt:      hiredAt,
    }
    if err := s.empRepo.Create(emp); err != nil {
        return nil, err
    }
    return emp, nil
}