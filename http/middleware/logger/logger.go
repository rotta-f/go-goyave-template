package logger

import (
	"time"

	"goyave.dev/goyave/v5"
)

type Middleware struct {
	goyave.Component
}

func (ctrl *Middleware) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		now := time.Now()

		next(response, request)

		ctrl.Logger().With(
			"method", request.Method(),
			"uri", request.Route.GetFullURI(),
			"status", response.GetStatus(),
			"duration", time.Since(now),
		).Info("Request received")
	}
}
