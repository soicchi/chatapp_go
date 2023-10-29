package router

import (
	"chatapp/internal/infrastructure/database"
	"chatapp/internal/interface/handler"
	"chatapp/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Router struct {
	AuthHandler handler.AuthHandler
}

func InitRouter(e *echo.Echo, db *gorm.DB) {
	userRepo := database.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	authHandler := handler.NewAuthHandler(userUseCase)

	router := newRouter(*authHandler)
	router.setUpRouter(e)
}

func newRouter(authHandler handler.AuthHandler) *Router {
	return &Router{
		AuthHandler: authHandler,
	}
}

func (r *Router) setUpRouter(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/signup", r.AuthHandler.SignUp)
	auth.POST("/signin", r.AuthHandler.SignIn)
}
