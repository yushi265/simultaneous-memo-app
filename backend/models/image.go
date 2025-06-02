package models

import (
	"time"

	"gorm.io/gorm"
)

// Image represents an uploaded image with metadata
type Image struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Filename      string    `json:"filename" gorm:"not null"`
	OriginalName  string    `json:"original_name"`
	Path          string    `json:"path" gorm:"not null"`
	ThumbnailPath string    `json:"thumbnail_path"`
	Size          int64     `json:"size"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	ContentType   string    `json:"content_type"`
	PageID        *uint     `json:"page_id" gorm:"index"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateImage creates a new image record
func CreateImage(db *gorm.DB, image *Image) error {
	return db.Create(image).Error
}

// GetImageByID retrieves an image by ID
func GetImageByID(db *gorm.DB, id uint) (*Image, error) {
	var image Image
	err := db.First(&image, id).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImageByFilename retrieves an image by filename
func GetImageByFilename(db *gorm.DB, filename string) (*Image, error) {
	var image Image
	err := db.Where("filename = ?", filename).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImagesByPageID retrieves all images associated with a page
func GetImagesByPageID(db *gorm.DB, pageID uint) ([]Image, error) {
	var images []Image
	err := db.Where("page_id = ?", pageID).Order("created_at DESC").Find(&images).Error
	return images, err
}

// UpdateImage updates an existing image
func UpdateImage(db *gorm.DB, id uint, updates map[string]interface{}) error {
	return db.Model(&Image{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteImage deletes an image
func DeleteImage(db *gorm.DB, id uint) error {
	return db.Delete(&Image{}, id).Error
}

// GetOrphanedImages retrieves images not associated with any page
func GetOrphanedImages(db *gorm.DB, olderThan time.Time) ([]Image, error) {
	var images []Image
	err := db.Where("page_id IS NULL AND created_at < ?", olderThan).Find(&images).Error
	return images, err
}