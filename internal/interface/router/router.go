package router

import (
	"chatapp/internal/interface/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	AuthHandler handler.AuthHandler
}

func NewRouter(authHandler handler.AuthHandler) *Router {
	return &Router{
		AuthHandler: authHandler,
	}
}

func (r *Router) NewRouter(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/signup", r.AuthHandler.SignUp)
	auth.POST("/signin", r.AuthHandler.SignIn)
}
