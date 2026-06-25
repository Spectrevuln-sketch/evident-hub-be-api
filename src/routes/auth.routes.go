package routes

import (
	"evidence-hub-be/src/core/utils"
	"evidence-hub-be/src/handler/auth"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	route   *gin.RouterGroup
	handler *auth.Handler
}

func NewAuthRoutes(group *gin.RouterGroup, handler *auth.Handler) *AuthRoutes {
	return &AuthRoutes{
		route:   group,
		handler: handler,
	}
}

func (r *AuthRoutes) SetupRoutes() {
	r.route.POST("/register", utils.Wrap(r.handler.SignUpHandler))
	r.route.POST("/login", utils.Wrap(r.handler.SigninHandler))
	r.route.GET("/token-check", utils.Wrap(r.handler.TokenCheck))
	r.route.POST("/change-password", utils.Wrap(r.handler.ChangePassword))
}
