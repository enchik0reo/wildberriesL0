package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/enchik0reo/wildberriesL0/internal/config"
	"github.com/enchik0reo/wildberriesL0/internal/handlers"
	"github.com/enchik0reo/wildberriesL0/internal/pkg"
	"github.com/enchik0reo/wildberriesL0/internal/repository"
	"github.com/enchik0reo/wildberriesL0/internal/repository/cache"
	"github.com/enchik0reo/wildberriesL0/internal/repository/storage/psql"
	"github.com/enchik0reo/wildberriesL0/internal/server"
)

type App struct {
	cfg        *config.Config
	httpServer *server.Server
	repo       *repository.Repository
	nats       *pkg.Stan
	handler    *handlers.Handler
}

func New() *App {
	var err error
	a := &App{}

	a.cfg, err = config.Load()
	if err != nil {
		log.Fatalf("error loading config variables: %s", err.Error())
	}

	a.httpServer = server.New()

	cache := cache.New()

	psql, err := psql.New(a.cfg.DB.Host, a.cfg.DB.Port, a.cfg.DB.User, a.cfg.DB.Password, a.cfg.DB.DBName, a.cfg.DB.SSLMode)
	if err != nil {
		log.Fatalf("failed to inicialise DB: %s", err.Error())
	}

	a.repo = repository.New(psql, cache)

	a.nats, err = pkg.NatsConnect(a.cfg.Nats.ClusterID, a.cfg.Nats.ClientID, a.cfg.Nats.URL)
	if err != nil {
		log.Fatalf("failed to connect nuts: %s", err.Error())
	}

	a.handler = handlers.New(a.repo)

	log.Println("The New Application Created")
	return a
}

func (a *App) Run() {
	ctx := context.Background()

	go func() {
		if err := a.nats.GetMsg(ctx, a.repo); err != nil {
			log.Fatalf("error occured while nuts got a message: %s", err.Error())
		}
	}()

	go func() {
		if err := a.httpServer.Run(a.cfg.Port, a.handler.InitRoute()); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Fatalf("error occured while http server started: %s", err.Error())
		}
	}()

	log.Print("The App Successfully Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("The App Shutting Down")

	if err := a.httpServer.Stop(ctx); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

	if err := a.nats.Conn.Close(); err != nil {
		log.Fatalf("error occured on nuts connection close: %s", err.Error())
	}

	if err := a.repo.Stop(ctx); err != nil {
		log.Fatalf("error occured on db connection close: %s", err.Error())
	}

	log.Print("The App Successfully Shutted Down")
}
