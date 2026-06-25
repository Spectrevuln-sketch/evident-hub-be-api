package routes

import (
	"evidence-hub-be/src/core/utils"
	"evidence-hub-be/src/handler/leaderboard"

	"github.com/gin-gonic/gin"
)

type LeaderBoardRoutes struct {
	route   *gin.RouterGroup
	handler *leaderboard.Handler
}

func NewLeaderBoardRoutes(group *gin.RouterGroup, handler *leaderboard.Handler) *LeaderBoardRoutes {
	return &LeaderBoardRoutes{
		route:   group,
		handler: handler,
	}
}

func (r *LeaderBoardRoutes) SetupRoutes() {
	leaderboard := r.route.Group("/leaderboard")
	leaderboard.GET("/", utils.Wrap(r.handler.GetAllLeaderBoard))
	leaderboard.GET("/:id", utils.Wrap(r.handler.GetLeaderBoardById))
	leaderboard.GET("/me", utils.Wrap(r.handler.GetMyLeaderBoard))

}
