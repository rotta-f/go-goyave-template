package extra

import (
	"goyave.dev/goyave/v5"
	"goyave.dev/template/database/model"
)

const (
	contextKeyParamsUser  = "params.user"
	contextKeyCurrentUser = "authentication.user"
)

func SetParamsUser(request *goyave.Request, user *model.User) {
	request.Extra[contextKeyParamsUser] = user
}

func GetParamsUser(request *goyave.Request) *model.User {
	user, ok := request.Extra[contextKeyParamsUser].(*model.User)
	if !ok {
		return nil
	}
	request.Context()
	return user
}

func SetCurrentUser(request *goyave.Request, user *model.User) {
	request.Extra[contextKeyCurrentUser] = user
}

func GetCurrentUser(request *goyave.Request) *model.User {
	user, ok := request.Extra[contextKeyCurrentUser].(*model.User)
	if !ok {
		return nil
	}
	return user
}
