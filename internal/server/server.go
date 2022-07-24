package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"restapi/internal/config"
	"restapi/internal/db/postgres"
	"restapi/internal/db/redis"
	"restapi/internal/logger"
	"syscall"
)

func Start() error {
	postgresClient, err := postgres.NewClient()
	if err != nil {
		return err
	}
	defer postgresClient.Close()

	redisClient, err := redis.NewClient()
	if err != nil {
		return err
	}
	defer redisClient.Close()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cfg().APPPort),
		Handler: NewRouter(postgresClient, redisClient),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			logger.Log().Err(err).Msg("failed to shutdown server")
		}
	}()

	logger.Log().Info().Msgf("starting server on %s", httpServer.Addr)
	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed

	logger.Log().Info().Msg("stopped server gracefully")
	return nil
}
