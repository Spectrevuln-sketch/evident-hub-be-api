package routes

import (
	"evidence-hub-be/src/core/utils"
	"evidence-hub-be/src/handler/users"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	route   *gin.RouterGroup
	handler *users.Handler
}

func NewUserRoutes(group *gin.RouterGroup, handler *users.Handler) *UserRoutes {
	return &UserRoutes{
		route:   group,
		handler: handler,
	}
}

func (r *UserRoutes) SetupRoutes() {
	usersRoute := r.route.Group("/users")
	usersRoute.GET("", utils.Wrap(r.handler.GetAllUsers))
}
