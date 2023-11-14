package auth

import (
	"fmt"
	"net/http"
	"strconv"

	"goyave.dev/goyave/v5"
	"goyave.dev/template/http/route/extra"
	"goyave.dev/template/service"
	"goyave.dev/template/service/user"
)

const (
	HeaderUserID = "X-User-ID"

	extraContextKey = "authentication.user"
)

type Middleware struct {
	goyave.Component
	UserService *user.Service
}

func (ctrl *Middleware) Init(server *goyave.Server) {
	ctrl.UserService = server.Service(service.User).(*user.Service)
	ctrl.Component.Init(server)
}

func (ctrl *Middleware) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		id, err := strconv.ParseInt(request.Header().Get(HeaderUserID), 10, 64)
		if err != nil || id <= 0 {
			ctrl.Logger().Error(fmt.Errorf("invalid user ID: %w", err))
			response.Status(http.StatusUnauthorized)
			return
		}

		user, err := ctrl.UserService.First(id)
		if err != nil {
			ctrl.Logger().Error(fmt.Errorf("user not found: %w", err))
			response.Status(http.StatusUnauthorized)
			return
		}

		extra.SetCurrentUser(request, user)

		next(response, request)
	}
}
