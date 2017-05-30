package main

import (
	"github.com/bradacina/inventory/logincookie"
	"github.com/bradacina/inventory/services"
)

func (app *app) GetCookieAuthentication() *logincookie.CookieAuthentication {
	return app.cookieAuth
}

func (app *app) GetUserService() services.UserServicer {
	return app.userService
}

func (app *app) GetLoginRoute() string {
	return RouteLogin
}
