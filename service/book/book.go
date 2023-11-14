package book

import (
	"github.com/go-faker/faker/v4"
	"goyave.dev/filter"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/template/database/model"
	"goyave.dev/template/service"
)

type Repository interface {
	Create(book *model.Book) error
	Paginate(request *filter.Request) (*database.Paginator[*model.Book], error)
}

type Service struct {
	repository Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) Name() string {
	return service.Book
}

func (s *Service) Create(owner *model.User) (*model.Book, error) {
	// This is the business logic of the service.
	book := &model.Book{
		Owner: *owner,
		Title: faker.Sentence(),
	}

	if err := s.repository.Create(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (s *Service) Paginate(request *filter.Request) (*database.Paginator[*model.Book], error) {
	return s.repository.Paginate(request)
}
