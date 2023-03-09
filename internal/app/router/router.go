package router

import (
	"github.com/antnzr/chat-go/internal/app/controller"
	"github.com/gin-gonic/gin"
)

type appRouter struct {
	controller controller.Controller
	engine     *gin.Engine
}

func NewAppRouter(engine *gin.Engine, controller controller.Controller) *appRouter {
	return &appRouter{
		engine:     engine,
		controller: controller,
	}
}

func (r *appRouter) Setup() {
	v1 := r.engine.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", r.controller.Auth.Signup)
			auth.POST("/login", r.controller.Auth.Login)
		}

		users := v1.Group("/users")
		{
			users.GET("/me", r.controller.User.GetMe)
		}
	}
}
