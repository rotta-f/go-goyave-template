package model

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `gorm:"type:text"`
	Email string `gorm:"type:text;uniqueIndex"`

	OwnedBooks []Book `gorm:"foreignKey:OwnerID"`
}

func UserGenerator() *User {
	user := &User{}
	user.Name = faker.Name()

	user.Email = faker.Email(
		options.WithGenerateUniqueValues(true),
	)
	return user
}
