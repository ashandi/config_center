package application

import (
	"config_center/internal/config"
	"config_center/internal/repositories"
	"config_center/internal/server"
	"config_center/internal/server/handlers/config_request_handler"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (app *App) Start(ctx context.Context, errChan chan<- error, gracefully GracefulQuitFn) error {
	db, err := sqlx.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		app.cfg.MysqlDbUser,
		app.cfg.MysqlDbPassword,
		app.cfg.MysqlDbHost,
		app.cfg.MysqlDbName,
	))
	if err != nil {
		return err
	}

	depsRepo := repositories.NewDependenciesRepository(db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     app.cfg.RedisHost,
		Password: app.cfg.RedisPassword,
		DB:       app.cfg.RedisDb,
	})

	srv := app.initServer(depsRepo, rdb)
	go gracefully(func() {
		srv.Run(ctx, errChan)
	})

	return nil
}

func (app *App) initServer(depsRepo *repositories.DependenciesRepository, rdb *redis.Client) *server.Server {
	srv := server.New(app.cfg, app.logger)
	srv.SetupHandlers(
		config_request_handler.NewBaseHandler(depsRepo, rdb, app.cfg, app.logger),
	)

	return srv
}
