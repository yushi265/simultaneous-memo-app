package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Image represents an uploaded image with metadata
type Image struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	WorkspaceID   uuid.UUID `json:"workspace_id" gorm:"type:uuid;not null"`
	Filename      string    `json:"filename" gorm:"not null"`
	OriginalName  string    `json:"original_name"`
	Path          string    `json:"path" gorm:"not null"`
	ThumbnailPath string    `json:"thumbnail_path"`
	Size          int64     `json:"size"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	ContentType   string    `json:"content_type"`
	PageID        *uuid.UUID `json:"page_id" gorm:"type:uuid;index"`
	UserID        uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Relationships
	Workspace     Workspace  `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Page          *Page      `gorm:"foreignKey:PageID" json:"page,omitempty"`
	User          User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CreateImage creates a new image record
func CreateImage(db *gorm.DB, image *Image) error {
	return db.Create(image).Error
}

// GetImageByID retrieves an image by ID
func GetImageByID(db *gorm.DB, id uint, workspaceID uuid.UUID) (*Image, error) {
	var image Image
	err := db.Where("id = ? AND workspace_id = ?", id, workspaceID).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImageByFilename retrieves an image by filename
func GetImageByFilename(db *gorm.DB, filename string, workspaceID uuid.UUID) (*Image, error) {
	var image Image
	err := db.Where("filename = ? AND workspace_id = ?", filename, workspaceID).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImagesByPageID retrieves all images associated with a page
func GetImagesByPageID(db *gorm.DB, pageID uuid.UUID, workspaceID uuid.UUID) ([]Image, error) {
	var images []Image
	err := db.Where("page_id = ? AND workspace_id = ?", pageID, workspaceID).Order("created_at DESC").Find(&images).Error
	return images, err
}

// UpdateImage updates an existing image
func UpdateImage(db *gorm.DB, id uint, workspaceID uuid.UUID, updates map[string]interface{}) error {
	return db.Model(&Image{}).Where("id = ? AND workspace_id = ?", id, workspaceID).Updates(updates).Error
}

// DeleteImage deletes an image
func DeleteImage(db *gorm.DB, id uint, workspaceID uuid.UUID) error {
	return db.Where("id = ? AND workspace_id = ?", id, workspaceID).Delete(&Image{}).Error
}

// GetOrphanedImages retrieves images not associated with any page in a workspace
func GetOrphanedImages(db *gorm.DB, workspaceID uuid.UUID, olderThan time.Time) ([]Image, error) {
	var images []Image
	err := db.Where("workspace_id = ? AND page_id IS NULL AND created_at < ?", workspaceID, olderThan).Find(&images).Error
	return images, err
}