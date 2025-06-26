package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	// Enable UUID extension for PostgreSQL
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	
	// Migrate models in the correct order
	return db.AutoMigrate(
		&User{},
		&Workspace{},
		&WorkspaceMember{},
		&WorkspaceInvitation{},
		&Page{},
		&Image{},
		&File{},
	)
}