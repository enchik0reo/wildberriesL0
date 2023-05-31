package app

import (
	"context"
	"errors"
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
	"github.com/enchik0reo/wildberriesL0/pkg/logging"
)

type App struct {
	log        *logging.Lgr
	cfg        *config.Config
	httpServer *server.Server
	service    *service.Service
	handler    *handlers.Handler
}

func New() *App {
	var err error
	a := &App{}

	a.log = logging.New()

	a.cfg, err = config.Load()
	if err != nil {
		a.log.Fatalf("failed to load config variables: %s", err.Error())
	}

	a.httpServer = server.New()

	cache := cache.New(a.cfg.CacheSize)

	psql, err := psql.New(a.cfg.DB.Host, a.cfg.DB.Port, a.cfg.DB.User, a.cfg.DB.Password, a.cfg.DB.DBName, a.cfg.DB.SSLMode, a.log)
	if err != nil {
		a.log.Fatalf("failed to initialize DB: %s", err.Error())
	}

	repo, err := repository.New(psql, cache)
	if err != nil {
		a.log.Fatalf("failed to create repository: %s", err.Error())
	}

	nats, err := nats.New(a.cfg.Nats.ClusterID, a.cfg.Nats.ClientID, a.cfg.Nats.URL, a.log)
	if err != nil {
		a.log.Fatalf("failed to connect nuts: %s", err.Error())
	}

	a.service = service.New(nats, repo, a.log)

	a.handler = handlers.New(repo)

	a.log.Infoln("New application created")
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
			a.log.Fatalf("error occured while http server started: %s", err.Error())
		}
	}()

	a.log.Info("Application successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	a.log.Info("Application is shutting down")

	if err := a.httpServer.Stop(ctx); err != nil {
		a.log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

	if err := a.service.Stop(ctx); err != nil {
		a.log.Fatalf("error occured on service shutting down: %s", err.Error())
	}

	a.log.Info("Application successfully shutted down")
}
