package routes

import (
	"evidence-hub-be/src/core/utils"
	"evidence-hub-be/src/handler/evident"

	"github.com/gin-gonic/gin"
)

type EvidentRoutes struct {
	route   *gin.RouterGroup
	handler *evident.Handler
}

func NewEvidentRoutes(group *gin.RouterGroup, handler *evident.Handler) *EvidentRoutes {
	return &EvidentRoutes{
		route:   group,
		handler: handler,
	}
}

func (r *EvidentRoutes) SetupRoutes() {
	evident := r.route.Group("/evidents")
	evident.GET("/", utils.Wrap(r.handler.GetAllEvident))
	evident.GET("/:id", utils.Wrap(r.handler.GetEvidentById))
	evident.POST("", utils.Wrap(r.handler.CreateEvidentHandler))
	evident.PATCH("/revisi/:id", utils.Wrap(r.handler.RevisiEvidentHandler))

}
