package server

import (
	"context"
	"github.com/edermanoel94/pismo/internal/infra/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func New() *echo.Echo {

	e := echo.New()

	e.Use(
	//middleware.Logger(),
	//middleware.Recover(),
	//middleware.CORS(),
	)

	e.GET("/health", healthCheck)

	return e
}

func Start(e *echo.Echo) {

	s := &http.Server{
		Addr:         config.Config().GetString("server.addr"),
		ReadTimeout:  config.Config().GetDuration("server.timeout.write-seconds") * time.Second,
		WriteTimeout: config.Config().GetDuration("server.timeout.read-seconds") * time.Second,
	}

	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
