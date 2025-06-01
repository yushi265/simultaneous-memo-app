package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Page struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Content   datatypes.JSON `json:"content" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type BlockContent struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Content  interface{}    `json:"content"`
	Children []BlockContent `json:"children,omitempty"`
}

// CreatePage creates a new page
func CreatePage(db *gorm.DB, page *Page) error {
	return db.Create(page).Error
}

// GetPageByID retrieves a page by ID
func GetPageByID(db *gorm.DB, id uint) (*Page, error) {
	var page Page
	err := db.First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// GetAllPages retrieves all pages
func GetAllPages(db *gorm.DB) ([]Page, error) {
	var pages []Page
	err := db.Order("updated_at DESC").Find(&pages).Error
	return pages, err
}

// UpdatePage updates an existing page
func UpdatePage(db *gorm.DB, id uint, updates map[string]interface{}) error {
	return db.Model(&Page{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePage deletes a page
func DeletePage(db *gorm.DB, id uint) error {
	return db.Delete(&Page{}, id).Error
}