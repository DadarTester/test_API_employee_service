package repository

import (
    "errors"
    "gorm.io/gorm"
    "org-api/internal/models"
    apperrors "org-api/internal/errors"
)

type EmployeeRepository struct {
    db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
    return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(emp *models.Employee) error {
    return r.db.Create(emp).Error
}

func (r *EmployeeRepository) GetByDepartment(deptID uint) ([]models.Employee, error) {
    var employees []models.Employee
    err := r.db.Where("department_id = ?", deptID).Order("created_at ASC").Find(&employees).Error
    return employees, err
}

func (r *EmployeeRepository) DeleteByDepartment(deptID uint) error {
    return r.db.Where("department_id = ?", deptID).Delete(&models.Employee{}).Error
}

func (r *EmployeeRepository) ReassignDepartment(oldDeptID, newDeptID uint) error {
    return r.db.Model(&models.Employee{}).Where("department_id = ?", oldDeptID).Update("department_id", newDeptID).Error
}

func (r *EmployeeRepository) GetByID(id uint) (*models.Employee, error) {
    var emp models.Employee
    err := r.db.First(&emp, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, apperrors.ErrEmployeeNotFound
    }
    return &emp, err
}