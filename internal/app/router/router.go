package router

import (
	"fmt"
	"net/http"

	"github.com/antnzr/chat-go/internal/app/controller"
	"github.com/gin-gonic/gin"
)

type appRouter struct {
	engine     *gin.Engine
	controller *controller.Controller
	auth       gin.HandlerFunc
}

func NewAppRouter(engine *gin.Engine, controller *controller.Controller, auth gin.HandlerFunc) *appRouter {
	return &appRouter{
		engine:     engine,
		controller: controller,
		auth:       auth,
	}
}

func (r *appRouter) Setup() {
	v1 := r.engine.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", r.controller.Auth.Signup)
			auth.POST("/login", r.controller.Auth.Login)
			auth.GET("/refresh", r.controller.Auth.Refresh)
			auth.GET("/logout", r.auth, r.controller.Auth.Logout)
		}

		users := v1.Group("/users")
		{
			users.GET("/me", r.auth, r.controller.User.GetMe)
		}
	}

	r.engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusTeapot)
	})

	r.engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Route %s not found", ctx.Request.URL),
		})
	})
}
