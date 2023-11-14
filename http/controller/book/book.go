package book

import (
	"errors"
	"net/http"

	"goyave.dev/filter"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/util/typeutil"
	"goyave.dev/template/database/model"
	"goyave.dev/template/http/dto"
	"goyave.dev/template/http/route/extra"
	"goyave.dev/template/service"
)

type Service interface {
	Create(owner *model.User) (*model.Book, error)
	Paginate(request *filter.Request) (*database.Paginator[*model.Book], error)
}

type Controller struct {
	goyave.Component
	BookService Service
}

func (ctrl *Controller) Init(server *goyave.Server) {
	ctrl.BookService = server.Service(service.Book).(Service)
	ctrl.Component.Init(server)
}

func (ctrl *Controller) RegisterRoutes(router *goyave.Router) {
	subrouter := router.Subrouter("/books")

	subrouter.Get("/", ctrl.Index).ValidateQuery(filter.Validation)
	subrouter.Post("/", ctrl.NewBook)
}

func (ctrl *Controller) NewBook(response *goyave.Response, request *goyave.Request) {
	user := extra.GetCurrentUser(request)
	if user == nil {
		ctrl.Logger().Error(errors.New("current user not found"))
		response.Status(http.StatusUnauthorized)
		return
	}

	targetUser := extra.GetParamsUser(request)
	if targetUser != nil && targetUser.ID != user.ID {
		ctrl.Logger().Error(errors.New("not allowed to create a book as someone else"))
		response.Status(http.StatusForbidden)
		return
	}

	book, err := ctrl.BookService.Create(user)
	if err != nil {
		ctrl.Logger().Error(err)
		response.Error(err)
		return
	}

	dto := typeutil.MustConvert[dto.Book](book)
	response.JSON(http.StatusCreated, dto)
}

func (ctrl *Controller) Index(response *goyave.Response, request *goyave.Request) {
	ctrl.Logger().Info("Indexing books")

	paginator, err := ctrl.BookService.Paginate(filter.NewRequest(request.Query))
	if response.WriteDBError(err) {
		return
	}

	// Convert to DTO and write response
	dto := typeutil.MustConvert[database.PaginatorDTO[dto.Book]](paginator)
	response.JSON(http.StatusOK, dto)
}
