package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"goyave.dev/filter"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/template/database/model"
	"goyave.dev/template/service/book"
)

type BookSQL struct {
	DB *gorm.DB
}

func NewBookSQL(db *gorm.DB) *BookSQL {
	return &BookSQL{
		DB: db,
	}
}

// Assert BookSQL implements book.Repository.
// Enforce dependency inversion.
// Service should not depend on the database.
var _ book.Repository = (*BookSQL)(nil)

func (r *BookSQL) Create(book *model.Book) error {
	book.ID = 0
	result := r.DB.
		Clauses(clause.Returning{}).
		Save(book)
	return result.Error
}

func (r *BookSQL) Paginate(request *filter.Request) (*database.Paginator[*model.Book], error) {
	settings := filter.Settings[*model.Book]{
		DefaultSort: []*filter.Sort{
			{Field: "title", Order: filter.SortAscending},
		},
	}

	books := []*model.Book{}
	paginator, tx := settings.Scope(r.DB, request, &books)

	return paginator, tx.Error
}
