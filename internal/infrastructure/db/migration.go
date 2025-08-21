package db

import (
	"gorm.io/gorm"
	"tasius.my.id/todolistapi/internal/domain/entities"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
	)
}
