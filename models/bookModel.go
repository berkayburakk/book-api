package models

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	ID          uint `json:"-" gorm:"primaryKey"`
	BookName    string
	Category    string
	Author      string
	Barcode     string
	CreatedDate time.Time      `json:"-"`
	UpdatedDate time.Time      `json:"-"`
	DeletedDate gorm.DeletedAt `json:"-"`
}
