package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/MindscapeHQ/raygun4go"
	"github.com/blacklane/warsaw/logger"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload" // Loads environment variables from .env file

	"github.com/blacklane/worker/config"
	"github.com/blacklane/worker/external"
	"github.com/blacklane/worker/internal/handlers"
	"github.com/blacklane/worker/internal/middleware"
)

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C
	signal.Notify(signalChan, os.Kill)      // a.k.a kill

	// when closed the program should exit
	idleConnsClosed := make(chan bool)

	cfg, err := config.Parse()
	if err != nil {
		logger.Error("config_error", err).Msg("Could not parse environment variables")
		panic(err.Error())
	}

	initKafkaListeners(&cfg)

	raygun := startRaygunErrorTracking(cfg)
	errorReporter := external.NewRaygun(raygun)
	if raygun != nil {
		defer raygun.HandleError()
	}
	errorReporter.SendError(errors.New("just use the error to avoid compile errors"))

	router := initRouter(cfg)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,

		IdleTimeout:       cfg.IdleTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	// handle graceful shutdown in another goroutine
	go gracefulShutdown(signalChan, idleConnsClosed, &cfg, server)

	logger.Event("server_start").Msg(fmt.Sprintf("starting server on %v", cfg.ServerPort))

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Error("server_error", err).Msg(fmt.Sprintf("Error starting server on %v", cfg.ServerPort))
	} else {
		logger.Event("server_shutdown").Msg("Shutting server down...")
	}

	<-idleConnsClosed
}

func startRaygunErrorTracking(cfg config.Config) *raygun4go.Client {
	if cfg.RaygunEnabled && cfg.RaygunAPIKey != "" {
		raygun, err := raygun4go.New(cfg.RaygunAppName, cfg.RaygunAPIKey)
		if err != nil {
			logger.Error("raygun_error", err).Msg("Error starting Raygun client")
			panic(err)
		}

		return raygun
	}

	logger.Event("warn").Msg("Raygun is not enabled or configured")
	return nil
}

func initRouter(cfg config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(router chi.Router) {
		// Routes that will be logged
		router.Use(middleware.NewLoggerMiddleware(cfg.AppName))

		// add handlers here
	})
	router.Get("/ping", handlers.PingHandler())
	return router
}

func initKafkaListeners(cfg *config.Config) {
	// init kafka consumers here
	//reader := consumer.GetKafkaReader(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaConsumerGroup)
	//consumer.StartKafkaConsumer(ridesReader, listener.SomeKafkaListener)
}

func gracefulShutdown(
	signalChan chan os.Signal,
	idleConnsClosed chan bool,
	cfg *config.Config,
	server *http.Server) {

	sig := <-signalChan
	logger.Event("server_shutdown").Msg(fmt.Sprintf("received signal: %q, starting graceful shutdown...", sig.String()))

	ctx, done := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer done() // avoid a context leak

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server_error", err).Msg("error during gracefully shutdown")
	}

	logger.Event("server_shutdown").Msg("Shutdown finished, goodbye.")

	close(idleConnsClosed)
}
