package user

import (
	"net/http"

	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/util/typeutil"
	"goyave.dev/template/database/model"
	"goyave.dev/template/http/controller/book"
	"goyave.dev/template/http/dto"
	"goyave.dev/template/http/route/extra"
	"goyave.dev/template/service"
	"goyave.dev/template/service/user"
)

type Service interface {
	First(id int64) (*model.User, error)
	Paginate(page int, pageSize int) (*database.Paginator[*model.User], error)
}

type Controller struct {
	goyave.Component
	UserService Service
}

func (ctrl *Controller) Init(server *goyave.Server) {
	ctrl.UserService = server.Service(service.User).(*user.Service)
	ctrl.Component.Init(server)
}

func (ctrl *Controller) RegisterRoutes(router *goyave.Router) {
	subrouter := router.Subrouter("/users")

	subrouter.Get("/", ctrl.Index).ValidateQuery(IndexRequest)

	byId := subrouter.Subrouter("/{userId:[0-9+]}")
	byId.Middleware(&Middleware{})
	byId.Get("/", ctrl.Show)

	byId.Controller(&book.Controller{})
}

func (ctrl *Controller) Index(response *goyave.Response, request *goyave.Request) {
	ctrl.Logger().Info("Indexing users")
	query := typeutil.MustConvert[dto.Index](request.Query)

	paginator, err := ctrl.UserService.Paginate(query.Page.Default(1), query.PerPage.Default(20))
	if response.WriteDBError(err) {
		return
	}

	// Convert to DTO and write response
	dto := typeutil.MustConvert[database.PaginatorDTO[dto.User]](paginator)
	response.JSON(http.StatusOK, dto)
}

func (ctrl *Controller) Show(response *goyave.Response, request *goyave.Request) {
	user := extra.GetParamsUser(request)

	response.JSON(http.StatusOK, user)
}
