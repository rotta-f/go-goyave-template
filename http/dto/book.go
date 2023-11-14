package dto

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Book struct {
	ID        uint      `json:"id"`
	Owner     User      `json:"owner"`
	Title     string    `json:"title"`
	Readers   []User    `json:"readers"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt null.Time `json:"updatedAt"`
	DeletedAt null.Time `json:"deletedAt"`
}
