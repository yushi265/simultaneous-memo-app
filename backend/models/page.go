package models

import (
	"encoding/json"
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

// ImageReference represents an image reference within page content
type ImageReference struct {
	ImageID uint   `json:"image_id"`
	Src     string `json:"src"`
	Alt     string `json:"alt,omitempty"`
}

// BlockContent represents a content block in the page
type BlockContent struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Content  interface{}    `json:"content"`
	Attrs    interface{}    `json:"attrs,omitempty"`     // For image blocks, contains ImageReference
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

// ExtractImageReferences extracts all image references from page content
func ExtractImageReferences(content datatypes.JSON) ([]uint, error) {
	var imageIDs []uint
	if content == nil {
		return imageIDs, nil
	}

	// Parse the JSON content
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	// Extract image IDs from the content structure
	extractFromBlock := func(block map[string]interface{}) {
		if blockType, ok := block["type"].(string); ok && blockType == "image" {
			if attrs, ok := block["attrs"].(map[string]interface{}); ok {
				if imageID, ok := attrs["image_id"].(float64); ok {
					imageIDs = append(imageIDs, uint(imageID))
				}
			}
		}
	}

	// Walk through the content structure
	var walkContent func(interface{})
	walkContent = func(content interface{}) {
		switch v := content.(type) {
		case map[string]interface{}:
			extractFromBlock(v)
			if children, ok := v["content"].([]interface{}); ok {
				for _, child := range children {
					walkContent(child)
				}
			}
		case []interface{}:
			for _, item := range v {
				walkContent(item)
			}
		}
	}

	if doc, ok := data["doc"].(map[string]interface{}); ok {
		walkContent(doc)
	} else if content, ok := data["content"]; ok {
		walkContent(content)
	}

	return imageIDs, nil
}

// UpdateImageReferences updates the page_id for images referenced in the content
func UpdateImageReferences(db *gorm.DB, pageID uint, content datatypes.JSON) error {
	imageIDs, err := ExtractImageReferences(content)
	if err != nil {
		return err
	}

	// First, unlink all images from this page
	if err := db.Model(&Image{}).Where("page_id = ?", pageID).Update("page_id", nil).Error; err != nil {
		return err
	}

	// Then, link the referenced images
	if len(imageIDs) > 0 {
		if err := db.Model(&Image{}).Where("id IN ?", imageIDs).Update("page_id", pageID).Error; err != nil {
			return err
		}
	}

	return nil
}