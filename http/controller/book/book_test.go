package book_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"goyave.dev/filter"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/config"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/util/testutil"
	"goyave.dev/template/database/model"
	"goyave.dev/template/http/controller/book"
	"goyave.dev/template/http/route/extra"
	"goyave.dev/template/service"
)

type BookServiceMock struct {
	createReturn *model.Book
	createError  error
}

func (s *BookServiceMock) Name() string {
	return service.Book
}

func (s *BookServiceMock) Create(owner *model.User) (*model.Book, error) {
	return s.createReturn, s.createError
}

func (s *BookServiceMock) Paginate(request *filter.Request) (*database.Paginator[*model.Book], error) {
	return nil, nil
}

func TestController_NewBook(t *testing.T) {
	User1 := &model.User{Model: gorm.Model{ID: 1}}
	User2 := &model.User{Model: gorm.Model{ID: 2}}

	testingSet := []struct {
		name           string
		userAuth       *model.User
		userParams     *model.User
		serviceMock    *BookServiceMock
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:           "success",
			userAuth:       User1,
			serviceMock:    &BookServiceMock{&model.Book{Model: gorm.Model{ID: 1}, Title: "test", Owner: *User1}, nil},
			expectedStatus: 201,
			expectedJSON: `{
				"id": 1,
				"createdAt": "0001-01-01T00:00:00Z",
				"updatedAt": "0001-01-01T00:00:00Z",
				"deletedAt": null,
				"owner": {
					"id": 1,
					"createdAt": "0001-01-01T00:00:00Z",
					"updatedAt": "0001-01-01T00:00:00Z",
					"deletedAt": null,
					"name": "",
					"email": ""
				},
				"title": "test",
				"readers": null
			}`,
		},
		{
			name:           "success_with_params",
			userAuth:       User2,
			userParams:     User2,
			serviceMock:    &BookServiceMock{&model.Book{}, nil},
			expectedStatus: 201,
		},
		{
			name:           "errors_with_params",
			userAuth:       User1,
			userParams:     User2,
			serviceMock:    &BookServiceMock{&model.Book{}, nil},
			expectedStatus: 403,
		},
		{
			name:           "service_error",
			userAuth:       User1,
			serviceMock:    &BookServiceMock{nil, errors.New("test")},
			expectedStatus: 500,
			expectedJSON: `{
				"error": "test"
			}`,
		},
	}
	for _, tt := range testingSet {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.NewTestServerWithOptions(nil, goyave.Options{Config: config.LoadDefault()}, nil)
			request := testutil.NewTestRequest("POST", "/books", nil)
			response, recorder := testutil.NewTestResponse(request)

			server.RegisterService(tt.serviceMock)

			extra.SetCurrentUser(request, tt.userAuth)
			extra.SetParamsUser(request, tt.userParams)

			SUT := book.Controller{}
			SUT.Init(server.Server)

			SUT.NewBook(response, request)

			assert.Equal(t, tt.expectedStatus, response.GetStatus())
			if tt.expectedJSON != "" {
				assert.JSONEq(t, tt.expectedJSON, recorder.Body.String())
			}
		})
	}
}
