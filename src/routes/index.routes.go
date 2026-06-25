package routes

import (
	"evidence-hub-be/src/core/config/envs"
	"evidence-hub-be/src/core/constants"
	"evidence-hub-be/src/core/middleware"
	"evidence-hub-be/src/core/pkg/token"
	"evidence-hub-be/src/handler/auth"
	"evidence-hub-be/src/handler/evident"
	healthcheck "evidence-hub-be/src/handler/healthCheck"
	"evidence-hub-be/src/handler/leaderboard"

	// "evidence-hub-be/src/handler/role"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	router     *gin.Engine
	cfg        *envs.Config
	middleware *middleware.Middleware
	token      *token.Token
}

func New(
	router *gin.Engine,
	cfg *envs.Config,
	token *token.Token,
) *Routes {

	mw := middleware.New(token)

	return &Routes{
		router:     router,
		cfg:        cfg,
		middleware: mw,
		token:      token,
	}
}

func (r *Routes) Init() {

	// Swagger
	r.router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// Static File
	r.router.Static("/uploads", constants.UploadDir)

	// Health Check
	r.router.GET("/", healthcheck.HealtCheckHandler)

	// V1 Group
	v1 := r.router.Group("/v1")

	// Auth
	authHandler := auth.NewHandler(
		r.cfg,
		r.token,
	)

	authRoutes := NewAuthRoutes(
		v1,
		authHandler,
	)

	authRoutes.SetupRoutes()
	// Evident
	evidentHandler := evident.NewHandler(
		r.cfg,
		r.token,
	)

	evidentRoutes := NewEvidentRoutes(
		v1,
		evidentHandler,
	)

	evidentRoutes.SetupRoutes()

	// LeaderBoard
	leaderBoardHandler := leaderboard.NewHandler(
		r.cfg,
		r.token,
	)

	leaderBoardRoutes := NewLeaderBoardRoutes(
		v1,
		leaderBoardHandler,
	)

	leaderBoardRoutes.SetupRoutes()

	// Role
	// roleHandler := role.NewHandler()

	// roleRoutes := NewRoleRoutes(
	// 	v1,
	// 	roleHandler,
	// 	r.middleware,
	// )

	// roleRoutes.SetupRoutes()

}
