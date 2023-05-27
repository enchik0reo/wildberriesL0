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
	"github.com/enchik0reo/wildberriesL0/internal/nats"
	"github.com/enchik0reo/wildberriesL0/internal/repository"
	"github.com/enchik0reo/wildberriesL0/internal/repository/cache"
	"github.com/enchik0reo/wildberriesL0/internal/repository/storage/psql"
	"github.com/enchik0reo/wildberriesL0/internal/server"
	"github.com/enchik0reo/wildberriesL0/internal/service"
)

type App struct {
	cfg        *config.Config
	httpServer *server.Server
	repo       *repository.Repository
	nats       *nats.Stan
	service    *service.Service
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
		log.Fatalf("failed to initialize DB: %s", err.Error())
	}

	a.repo = repository.New(psql, cache)

	a.nats, err = nats.New(a.cfg.Nats.ClusterID, a.cfg.Nats.ClientID, a.cfg.Nats.URL)
	if err != nil {
		log.Fatalf("failed to connect nuts: %s", err.Error())
	}

	a.service = service.New(a.nats, a.repo)

	a.handler = handlers.New(a.repo)

	log.Println("New application created")
	return a
}

func (a *App) Run() {
	ctx := context.Background()

	go a.service.Work(ctx)

	go func() {
		if err := a.httpServer.Run(a.cfg.Port, a.handler.InitRoute()); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Fatalf("error occured while http server started: %s", err.Error())
		}
	}()

	log.Print("Application successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Application is shutting down")

	if err := a.httpServer.Stop(ctx); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

	if err := a.nats.CloseConnect(); err != nil {
		log.Fatalf("error occured on nuts connection close: %s", err.Error())
	}

	if err := a.repo.CloseConnect(ctx); err != nil {
		log.Fatalf("error occured on db connection close: %s", err.Error())
	}

	log.Print("Application successfully shutted down")
}
