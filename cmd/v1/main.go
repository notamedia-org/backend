package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	docs "github.com/notamedia-org/backend/docs"
	"github.com/notamedia-org/backend/internal/api/login"
	"github.com/notamedia-org/backend/internal/api/register"
	_ "github.com/notamedia-org/backend/internal/api/user"
	"github.com/notamedia-org/backend/internal/config"
	"github.com/notamedia-org/backend/internal/database"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//nolint:gochecknoglobals
var (
	version     = "X.X.X"
	name        = "beats"
	description = "core beats"
)

type ApplicationInfo struct {
	version     string
	name        string
	description string
}

func main() {
	contentSwagger, err := os.ReadFile("./data/config-swagger.json")
	if err != nil {
		log.Fatalf("Cannot read file: %v\n", err)
	}

	swaggerConfig, err := config.ReadFromFileSwagger(contentSwagger)
	if err != nil {
		log.Fatalf("Cannot read config from file: %v\n", err)
	}

	appConfig, err := config.ReadFromEnv()
	if err != nil {
		log.Fatalf("Cannot read config from env: %v\n", err)
	}

	applicationInfo := ApplicationInfo{version, name, description}

	err = startServer(appConfig, applicationInfo, swaggerConfig)

	if err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}

	log.Println("Server shut down successfully")
}

func startServer(
	config *config.Config,
	info ApplicationInfo,
	swaggerConfig config.SwaggerConfig,
) error {
	r := mux.NewRouter()

	docs.SwaggerInfo.Title = swaggerConfig.Title
	docs.SwaggerInfo.Description = swaggerConfig.Description
	docs.SwaggerInfo.Version = swaggerConfig.Version
	docs.SwaggerInfo.Host = swaggerConfig.Host
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	db, err := database.StartPolling(config)
	if err != nil {
		return err
	}

	api := r.PathPrefix("/api/v1").Subrouter()
	api.PathPrefix(swaggerConfig.BasePath).Handler(httpSwagger.WrapHandler)
	api.HandleFunc("/user/login", login.Login(db, config)).Methods(http.MethodPost)
	api.HandleFunc("/user/register", register.Register(db, config)).Methods(http.MethodPost)

	const defaultTimeout = 30 * time.Second
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(api)

	server := http.Server{ //nolint:exhaustruct
		Addr:         fmt.Sprintf(":%v", config.Port),
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
		Handler:      handler,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	errs := make(chan error)

	go func() {
		log.Printf("Starting HTTP server on %s", server.Addr)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errs <- fmt.Errorf("could not start HTTP server on %s: %w", server.Addr, err)
		}
	}()

	select {
	// Wait for a signal to shutdown the server
	case sig := <-interrupt:
		log.Printf("Received signal: %q\n", sig)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		// Shutdown the server gracefully
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not gracefully stop HTTP server: %w", err)
		}

	// Server could not start at all
	case err := <-errs:
		return err
	}

	return nil
}
