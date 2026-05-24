package service

import (
    "strings"
    "org-api/internal/models"
    "org-api/internal/repository"
    apperrors "org-api/internal/errors"
    "gorm.io/gorm"
)

type DepartmentService struct {
    deptRepo *repository.DepartmentRepository
    empRepo  *repository.EmployeeRepository
    db       *gorm.DB
    maxDepth int
}

func NewDepartmentService(deptRepo *repository.DepartmentRepository, empRepo *repository.EmployeeRepository, db *gorm.DB, maxDepth int) *DepartmentService {
    return &DepartmentService{
        deptRepo: deptRepo,
        empRepo:  empRepo,
        db:       db,
        maxDepth: maxDepth,
    }
}

func (s *DepartmentService) Create(name string, parentID *uint) (*models.Department, error) {
    name = strings.TrimSpace(name)
    if name == "" || len(name) > 200 {
        return nil, apperrors.ErrInvalidInput
    }
    // уникальность имени в пределах одного родителя
    exists, err := s.deptRepo.ExistsByNameUnderParent(name, parentID, 0)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, apperrors.ErrDuplicateName
    }
    dept := &models.Department{
        Name:     name,
        ParentID: parentID,
    }
    if err := s.deptRepo.Create(dept); err != nil {
        return nil, err
    }
    return dept, nil
}

func (s *DepartmentService) GetByID(id uint, depth int, includeEmployees bool) (*models.Department, error) {
    if depth > s.maxDepth {
        depth = s.maxDepth
    }
    dept, err := s.deptRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    if includeEmployees {
        employees, err := s.empRepo.GetByDepartment(id)
        if err != nil {
            return nil, err
        }
        dept.Employees = employees
    }
    if depth > 0 {
        children, err := s.getChildrenTree(id, depth-1, includeEmployees)
        if err != nil {
            return nil, err
        }
        dept.Children = children
    }
    return dept, nil
}

func (s *DepartmentService) getChildrenTree(parentID uint, depth int, includeEmployees bool) ([]models.Department, error) {
    children, err := s.deptRepo.GetChildren(parentID)
    if err != nil {
        return nil, err
    }
    if depth == 0 {
        
        if includeEmployees {
            for i := range children {
                emps, _ := s.empRepo.GetByDepartment(children[i].ID)
                children[i].Employees = emps
            }
        }
        return children, nil
    }
    for i := range children {
        if includeEmployees {
            emps, _ := s.empRepo.GetByDepartment(children[i].ID)
            children[i].Employees = emps
        }
        subChildren, err := s.getChildrenTree(children[i].ID, depth-1, includeEmployees)
        if err != nil {
            return nil, err
        }
        children[i].Children = subChildren
    }
    return children, nil
}

func (s *DepartmentService) Update(id uint, name *string, parentID *uint) (*models.Department, error) {
    dept, err := s.deptRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    if name != nil {
        newName := strings.TrimSpace(*name)
        if newName == "" || len(newName) > 200 {
            return nil, apperrors.ErrInvalidInput
        }
      
        targetParent := dept.ParentID
        if parentID != nil {
            targetParent = parentID
        }
        exists, err := s.deptRepo.ExistsByNameUnderParent(newName, targetParent, id)
        if err != nil {
            return nil, err
        }
        if exists {
            return nil, apperrors.ErrDuplicateName
        }
        dept.Name = newName
    }
    if parentID != nil {
      
        if parentID != nil {
            return nil, apperrors.ErrCycleDetected
        }

        if parentID != nil {
          
            if _, err := s.deptRepo.GetByID(*parentID); err != nil {
                return nil, apperrors.ErrInvalidParent
            }
        
            if s.isDescendant(*parentID, id) {
                return nil, apperrors.ErrCycleDetected
            }
        }
        dept.ParentID = parentID
    }
    if err := s.deptRepo.Update(dept); err != nil {
        return nil, err
    }
    return dept, nil
}

func (s *DepartmentService) isDescendant(node, ancestor uint) bool {
    for {
        dept, err := s.deptRepo.GetByID(node)
        if err != nil || dept.ParentID == nil {
            return false
        }
        if *dept.ParentID == ancestor {
            return true
        }
        node = *dept.ParentID
    }
}

func (s *DepartmentService) Delete(id uint, mode string, reassignTo *uint) error {
   
    subtreeIDs, err := s.deptRepo.GetSubtreeIDs(id)
    if err != nil {
        return err
    }
    return s.db.Transaction(func(tx *gorm.DB) error {
        if mode == "cascade" {
          
            for _, did := range subtreeIDs {
                if err := s.empRepo.DeleteByDepartment(did); err != nil {
                    return err
                }
            }
            
            if err := tx.Where("id IN ?", subtreeIDs).Delete(&models.Department{}).Error; err != nil {
                return err
            }
        } else if mode == "reassign" {
            if reassignTo == nil {
                return apperrors.ErrInvalidInput
            }
           
            if _, err := s.deptRepo.GetByID(*reassignTo); err != nil {
                return err
            }
         
            for _, did := range subtreeIDs {
                if err := s.empRepo.ReassignDepartment(did, *reassignTo); err != nil {
                    return err
                }
            }
           
            if err := tx.Model(&models.Department{}).Where("parent_id = ?", id).Update("parent_id", reassignTo).Error; err != nil {
                return err
            }
       
            if err := tx.Delete(&models.Department{}, id).Error; err != nil {
                return err
            }
        } else {
            return apperrors.ErrInvalidInput
        }
        return nil
    })
}