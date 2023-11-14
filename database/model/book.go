package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title   string `gorm:"type:text"`
	Owner   User
	OwnerID uint
	Readers []User `gorm:"many2many:user_book_readers;"`
}
