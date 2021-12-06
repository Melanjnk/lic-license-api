package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/ozonmp/lic-license-api/internal/config"
	"github.com/ozonmp/lic-license-api/internal/database"
	"github.com/ozonmp/lic-license-api/internal/metrics"
	"github.com/ozonmp/lic-license-api/internal/pkg/logger"
	"github.com/ozonmp/lic-license-api/internal/server"
	"github.com/ozonmp/lic-license-api/internal/tracer"
)

var (
	batchSize uint = 2
)

func main() {
	ctx := context.Background()

	if err := config.ReadConfigYML("config.yml"); err != nil {
		logger.FatalKV(ctx, "Failed init configuration", "err", err)
	}
	cfg := config.GetConfigInstance()

	syncLogger := logger.InitLogger(ctx, cfg.Project.Debug, "service", cfg.Project.Name)
	defer syncLogger()

	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

	logger.InfoKV(ctx, fmt.Sprintf("Starting service: %s", cfg.Project.Name),
		"version", cfg.Project.Version,
		"commitHash", cfg.Project.CommitHash,
		"debug", cfg.Project.Debug,
		"environment", cfg.Project.Environment,
	)

	metrics.InitMetrics(&cfg)

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := database.NewPostgres(ctx, dsn, cfg.Database.Driver)
	if err != nil {
		logger.FatalKV(ctx, "Failed init postgres", "err", err)
	}
	defer func() {
		if errCl := db.Close(); errCl != nil {
			logger.ErrorKV(ctx, "failed close DB connection", "err", errCl)
		}
	}()

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			logger.ErrorKV(ctx, "Migration failed", "err", err)

			return
		}
	}

	tracing, err := tracer.NewTracer(ctx, cfg.Jaeger.Service, cfg.Jaeger.Host, cfg.Jaeger.Port)
	if err != nil {
		logger.ErrorKV(ctx, "Failed init tracing", "err", err)

		return
	}
	defer func() {
		if err := tracing.Close(); err != nil {
			logger.ErrorKV(ctx, "Failed close tracer", "err", err)
		}
	}()

	if err := server.NewGrpcServer(db, batchSize).Start(ctx, &cfg); err != nil {
		logger.ErrorKV(ctx, "Failed creating gRPC server", "err", err)

		return
	}
}
