package route

import (
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/cors"
	"goyave.dev/goyave/v5/middleware/parse"
	"goyave.dev/template/http/controller/book"
	"goyave.dev/template/http/controller/user"
	"goyave.dev/template/http/middleware/auth"
	"goyave.dev/template/http/middleware/logger"
)

func Register(server *goyave.Server, router *goyave.Router) {
	router.GlobalMiddleware(&logger.Middleware{})
	router.CORS(cors.Default())

	router.GlobalMiddleware(&parse.Middleware{})
	router.GlobalMiddleware(&auth.Middleware{})

	router.Controller(&user.Controller{})
	router.Controller(&book.Controller{})
}
