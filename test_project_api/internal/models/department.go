package models

import (
    "time"
)

type Department struct {
    ID        uint       `gorm:"primaryKey" json:"id"`
    Name      string     `gorm:"not null;size:200" json:"name"`
    ParentID  *uint      `json:"parent_id"`
    CreatedAt time.Time  `json:"created_at"`
    Children  []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
    Employees []Employee   `json:"employees,omitempty"`
}

func (Department) TableName() string {
    return "departments"
}