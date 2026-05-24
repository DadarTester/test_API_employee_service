package repository

import (
    "errors"
    "gorm.io/gorm"
    "org-api/internal/models"
    apperrors "org-api/internal/errors"
)

type DepartmentRepository struct {
    db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
    return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(dept *models.Department) error {
    return r.db.Create(dept).Error
}

func (r *DepartmentRepository) GetByID(id uint) (*models.Department, error) {
    var dept models.Department
    err := r.db.First(&dept, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, apperrors.ErrDepartmentNotFound
    }
    return &dept, err
}

func (r *DepartmentRepository) Update(dept *models.Department) error {
    return r.db.Save(dept).Error
}

func (r *DepartmentRepository) Delete(id uint) error {
    return r.db.Delete(&models.Department{}, id).Error
}

func (r *DepartmentRepository) GetChildren(parentID uint) ([]models.Department, error) {
    var children []models.Department
    err := r.db.Where("parent_id = ?", parentID).Find(&children).Error
    return children, err
}

func (r *DepartmentRepository) GetSubtreeIDs(rootID uint) ([]uint, error) {
    var ids []uint
    // рекурсивный запрос для PostgreSQL
    query := `
        WITH RECURSIVE subtree AS (
            SELECT id FROM departments WHERE id = ?
            UNION ALL
            SELECT d.id FROM departments d
            INNER JOIN subtree s ON d.parent_id = s.id
        )
        SELECT id FROM subtree;
    `
    err := r.db.Raw(query, rootID).Scan(&ids).Error
    return ids, err
}

func (r *DepartmentRepository) ExistsByNameUnderParent(name string, parentID *uint, excludeID uint) (bool, error) {
    var count int64
    query := r.db.Model(&models.Department{}).Where("name = ?", name)
    if parentID == nil {
        query = query.Where("parent_id IS NULL")
    } else {
        query = query.Where("parent_id = ?", *parentID)
    }
    if excludeID != 0 {
        query = query.Where("id != ?", excludeID)
    }
    err := query.Count(&count).Error
    return count > 0, err
}

func (r *DepartmentRepository) UpdateParent(id uint, newParentID *uint) error {
    return r.db.Model(&models.Department{}).Where("id = ?", id).Update("parent_id", newParentID).Error
}