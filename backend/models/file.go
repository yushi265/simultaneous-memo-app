package models

import (
	"time"
)

type File struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Filename     string    `gorm:"not null" json:"filename"`
	OriginalName string    `gorm:"not null" json:"original_name"`
	ContentType  string    `gorm:"not null" json:"content_type"`
	Size         int64     `gorm:"not null" json:"size"`
	Path         string    `gorm:"not null;unique" json:"path"`
	PageID       *uint     `json:"page_id,omitempty"`
	Page         *Page     `json:"page,omitempty" gorm:"constraint:OnDelete:SET NULL;"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// FileMetadata represents file information for API responses
type FileMetadata struct {
	ID           uint      `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"`
	URL          string    `json:"url"`
	PageID       *uint     `json:"page_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// ToMetadata converts File to FileMetadata with URL
func (f *File) ToMetadata(baseURL string) FileMetadata {
	return FileMetadata{
		ID:           f.ID,
		Filename:     f.Filename,
		OriginalName: f.OriginalName,
		ContentType:  f.ContentType,
		Size:         f.Size,
		URL:          baseURL + "/api/file/" + f.Filename,
		PageID:       f.PageID,
		CreatedAt:    f.CreatedAt,
	}
}