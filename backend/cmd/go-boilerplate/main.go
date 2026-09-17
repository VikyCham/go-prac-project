package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/VikyCham/go-boilerplate/internal/config"
	"github.com/VikyCham/go-boilerplate/internal/database"
	"github.com/VikyCham/go-boilerplate/internal/handler"
	"github.com/VikyCham/go-boilerplate/internal/logger"
	"github.com/VikyCham/go-boilerplate/internal/repository"
	"github.com/VikyCham/go-boilerplate/internal/router"
	"github.com/VikyCham/go-boilerplate/internal/server"
	"github.com/VikyCham/go-boilerplate/internal/service"
)

const DefaultContextTimeout = 30

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Initialize New Relic logger Service
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	if cfg.Primary.Env != "local" {
		if err := database.Migrate(context.Background(), &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("Failed to migrate database")
		}
	}

	// Initalize server
	srv, err := server.New(cfg, &log, loggerService)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize server")
	}

	// Initialize repositories, services, and handlers
	repos := repository.NewRepositories(srv)
	services, serviceErr := service.NewServices(srv, repos)
	if serviceErr != nil {
		log.Fatal().Err(serviceErr).Msg("Could not create services")
	}

	handlers := handler.NewHandlers(srv, services)

	// Initialize router
	r := router.NewRouter(srv, handlers, services)

	// Setup HTTP server
	srv.SetupHttpServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Start server
	go func ()  {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.ShutDown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	stop()
	cancel()

	log.Info().Msg("server exited properly")
}