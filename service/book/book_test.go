package book_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"goyave.dev/filter"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/template/database/model"
	"goyave.dev/template/service/book"
)

type RepositoryMock struct {
	createError   error
	createID      uint
	createNbCalls int
}

func (r *RepositoryMock) Create(book *model.Book) error {
	book.ID = r.createID
	r.createNbCalls++
	return r.createError
}

func (r *RepositoryMock) Paginate(request *filter.Request) (*database.Paginator[*model.Book], error) {
	return nil, nil
}

func TestService_Create(t *testing.T) {
	User1 := model.User{Model: gorm.Model{ID: 1}}

	// Easier to test business logic without
	// testing sets.
	// We can go deeper in testing the business logic.

	t.Run("success", func(t *testing.T) {
		repositoryMock := &RepositoryMock{
			createError: nil,
			createID:    1,
		}
		SUT := book.NewService(repositoryMock)
		book, err := SUT.Create(&User1)

		require.NoError(t, err)
		require.NotNil(t, book)
		assert.Equal(t, 1, repositoryMock.createNbCalls)
		assert.Equal(t, uint(1), book.ID)
		assert.Equal(t, User1, book.Owner)
		assert.NotEmpty(t, book.Title)
	})

	t.Run("error", func(t *testing.T) {
		repositoryMock := &RepositoryMock{
			createError:   gorm.ErrInvalidData,
			createID:      0,
			createNbCalls: 0,
		}
		SUT := book.NewService(repositoryMock)
		book, err := SUT.Create(&User1)

		require.Error(t, err)
		require.Nil(t, book)
		assert.Equal(t, 1, repositoryMock.createNbCalls)
		assert.Equal(t, gorm.ErrInvalidData, err)
	})
}
