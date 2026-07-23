package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/database"
	"github.com/Jeudry/adventist-stack/pkg/logger"
	members "github.com/Jeudry/adventist-stack/services/members"
	membersgrpc "github.com/Jeudry/adventist-stack/services/members/internal/grpc"
	"github.com/Jeudry/adventist-stack/services/members/internal/repository"
	"github.com/Jeudry/adventist-stack/services/members/internal/service"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	GRPCPort string `env:"MEMBERS_GRPC_PORT" envDefault:"50052"`
	Postgres config.Postgres
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("members", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.Postgres.AutoMigrate {
		if err := database.Migrate(cfg.Postgres.DSN, "members_schema_migrations", members.MigrationsFS, "migrations"); err != nil {
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

	repo := repository.NewMemberRepository(pool)
	svc := service.NewMemberService(repo)

	grpcServer := grpc.NewServer()
	membersv1.RegisterMemberServiceServer(grpcServer, membersgrpc.NewServer(svc))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Error("failed to listen", "port", cfg.GRPCPort, "err", err)
		os.Exit(1)
	}

	go func() {
		log.Info("members service listening", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down members service...")
	grpcServer.GracefulStop()
}
