package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	prayersv1 "github.com/Jeudry/adventist-stack/gen/prayers/v1"
	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/database"
	"github.com/Jeudry/adventist-stack/pkg/logger"
	"github.com/Jeudry/adventist-stack/services/prayers"
	prayersgrpc "github.com/Jeudry/adventist-stack/services/prayers/internal/grpc"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/repository"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	GRPCPort string `env:"PRAYERS_GRPC_PORT" envDefault:"50055"`
	Postgres config.Postgres
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("prayers", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.Postgres.AutoMigrate {
		if err := database.Migrate(cfg.Postgres.DSN, "prayers_schema_migrations", prayers.MigrationsFS, "migrations"); err != nil {
			log.Error("migrations failed", "err", err)
			os.Exit(1)
		}
		log.Info("migrations applied")
	}

	pool, err := database.Connect(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := repository.NewPrayerRepository(pool)
	svc := service.NewPrayerService(repo)

	grpcServer := grpc.NewServer()
	prayersv1.RegisterPrayerServiceServer(grpcServer, prayersgrpc.NewServer(svc))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Error("failed to listen", "port", cfg.GRPCPort, "err", err)
		os.Exit(1)
	}

	go func() {
		log.Info("prayers service listening", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down prayers service...")
	grpcServer.GracefulStop()
}
