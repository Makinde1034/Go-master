package models

import "gorm.io/gorm"

type Book struct {
	ID          uint    `gorm:"primary key;autoIncreament" json:"id"`
	Authur      *string `json:"authur"`
	Title       *string `json:"title"`
	Publisher   *string `json:"publisher"`
	Description *string `json:"description"`
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}
