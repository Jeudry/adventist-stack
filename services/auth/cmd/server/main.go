package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/Jeudry/adventist-stack/gen/auth/v1"
	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/database"
	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/pkg/logger"
	auth "github.com/Jeudry/adventist-stack/services/auth"
	authgrpc "github.com/Jeudry/adventist-stack/services/auth/internal/grpc"
	"github.com/Jeudry/adventist-stack/services/auth/internal/repository"
	"github.com/Jeudry/adventist-stack/services/auth/internal/service"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	GRPCPort string `env:"AUTH_GRPC_PORT" envDefault:"50051"`
	Postgres config.Postgres
	JWT      config.JWT
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("auth", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.Postgres.AutoMigrate {
		if err := database.Migrate(cfg.Postgres.DSN, "auth_schema_migrations", auth.MigrationsFS, "migrations"); err != nil {
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

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	repo := repository.NewUserRepository(pool)
	svc := service.New(repo, jwtManager)

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authgrpc.NewServer(svc))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Error("failed to listen", "port", cfg.GRPCPort, "err", err)
		os.Exit(1)
	}

	go func() {
		log.Info("servicio auth escuchando", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down auth service...")
	grpcServer.GracefulStop()
}
