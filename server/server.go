package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/config/envs"
	"evidence-hub-be/src/core/pkg/token"
	"evidence-hub-be/src/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var InterruptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

type Server struct {
	app        *gin.Engine
	httpServer *http.Server
	route      *routes.Routes
	cfg        *envs.Config
}

func New(cfg *envs.Config) *Server {

	config.ConnectDB(cfg)

	tokenManager := token.New(cfg)

	app := gin.New()

	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	route := routes.New(
		app,
		cfg,
		tokenManager,
	)

	return &Server{
		app:   app,
		route: route,
		cfg:   cfg,
	}
}

func (s *Server) Run() error {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		InterruptSignal...,
	)
	defer stop()

	waitGroup, ctx := errgroup.WithContext(ctx)

	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	s.route.Init()

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.AppConfig.Port),
		Handler: s.app,
	}

	waitGroup.Go(func() error {

		log.Println("Starting Project ....")
		log.Printf("HTTP Server Running On Port %s", s.cfg.AppConfig.Port)

		if err := s.httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {

		<-ctx.Done()

		log.Println("HTTP server shutdown gracefully")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		return s.httpServer.Shutdown(shutdownCtx)
	})

	return waitGroup.Wait()
}
