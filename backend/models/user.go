package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email        string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string      `gorm:"type:varchar(255);not null" json:"-"`
	Name         string      `gorm:"type:varchar(255)" json:"name"`
	AvatarURL    string      `gorm:"type:varchar(500)" json:"avatar_url"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	
	// Relationships
	Workspaces []WorkspaceMember `gorm:"foreignKey:UserID" json:"-"`
}

type Workspace struct {
	ID          uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string      `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description string      `gorm:"type:text" json:"description"`
	IsPersonal  bool        `gorm:"default:false" json:"is_personal"`
	OwnerID     uuid.UUID   `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	
	// Relationships
	Owner   User                `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members []WorkspaceMember   `gorm:"foreignKey:WorkspaceID" json:"members,omitempty"`
	Pages   []Page              `gorm:"foreignKey:WorkspaceID" json:"pages,omitempty"`
}

type WorkspaceMember struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	WorkspaceID  uuid.UUID   `gorm:"type:uuid;not null" json:"workspace_id"`
	UserID       uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	Role         string      `gorm:"type:varchar(50);not null;default:'member'" json:"role"` // owner, admin, member, viewer
	JoinedAt     time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`
	
	// Relationships
	Workspace Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Ensure unique constraint on workspace_id and user_id
func (WorkspaceMember) TableName() string {
	return "workspace_members"
}

// Hook to create personal workspace after user creation
func (u *User) AfterCreate(tx *gorm.DB) error {
	// Create personal workspace
	workspace := Workspace{
		Name:       u.Name + "'s Workspace",
		Slug:       u.ID.String(), // Use user ID as slug for personal workspace
		IsPersonal: true,
		OwnerID:    u.ID,
	}
	
	if err := tx.Create(&workspace).Error; err != nil {
		return err
	}
	
	// Add user as owner of the workspace
	member := WorkspaceMember{
		WorkspaceID: workspace.ID,
		UserID:      u.ID,
		Role:        "owner",
	}
	
	return tx.Create(&member).Error
}