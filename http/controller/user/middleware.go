package user

import (
	"net/http"
	"strconv"

	"goyave.dev/goyave/v5"
	"goyave.dev/template/http/route/extra"
	"goyave.dev/template/service"
)

type Middleware struct {
	goyave.Component
	UserService Service
}

func (ctrl *Middleware) Init(server *goyave.Server) {
	ctrl.UserService = server.Service(service.User).(Service)
	ctrl.Component.Init(server)
}

func (ctrl *Middleware) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		ctrl.Logger().With("params", request.RouteParams).Info("Loading user")
		userID, err := strconv.ParseInt(request.RouteParams["userId"], 10, 64)
		if err != nil {
			response.Status(http.StatusNotFound)
			return
		}

		ctrl.Logger().With("userId", userID).Info("Loading user")
		user, err := ctrl.UserService.First(userID)
		if response.WriteDBError(err) {
			return
		}

		extra.SetParamsUser(request, user)
		next(response, request)
	}
}
