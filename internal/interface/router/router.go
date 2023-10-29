package router

import (
	"chatapp/internal/infrastructure/database"
	"chatapp/internal/interface/handler"
	"chatapp/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Handlers struct {
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
}

func InitRouter(db *gorm.DB) *Handlers {
	userRepo := database.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	authHandler := handler.NewAuthHandler(userUseCase)

	userHandler := handler.NewUserHandler(userUseCase)

	handlers := &Handlers{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}

	return handlers
}

func (h *Handlers) SetUpRouter(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/signup", h.AuthHandler.SignUp)
	auth.POST("/signin", h.AuthHandler.SignIn)

	users := v1.Group("/users")
	users.GET("/:id", h.UserHandler.RetrieveUser)
	users.GET("/", h.UserHandler.ListUsers)
	users.PUT("/:id", h.UserHandler.UpdateUserInfo)
	users.DELETE("/:id", h.UserHandler.DeleteUser)
}
