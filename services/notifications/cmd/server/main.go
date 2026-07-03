// Command server arranca el microservicio de notificaciones (servidor gRPC).
package main

import (
	"context"
	"html/template"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	notificationsv1 "github.com/Jeudry/adventist-stack/gen/notifications/v1"
	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/logger"
	"github.com/Jeudry/adventist-stack/pkg/mailer"
	"github.com/Jeudry/adventist-stack/pkg/redis"
	notifications "github.com/Jeudry/adventist-stack/services/notifications"
	notifgrpc "github.com/Jeudry/adventist-stack/services/notifications/internal/grpc"
	"github.com/Jeudry/adventist-stack/services/notifications/internal/service"
)

// Config del servicio de notificaciones.
type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	GRPCPort string `env:"NOTIFICATIONS_GRPC_PORT" envDefault:"50052"`
	Redis    config.Redis
	SMTP     config.SMTP
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("notifications", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Redis (cache + pub/sub).
	rdb, err := redis.Connect(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Error("no se pudo conectar a redis", "err", err)
		os.Exit(1)
	}
	defer rdb.Close()

	// Plantillas de correo embebidas.
	tmpl, err := template.ParseFS(notifications.TemplatesFS, "templates/*.html")
	if err != nil {
		log.Error("no se pudieron parsear las plantillas", "err", err)
		os.Exit(1)
	}
	mail := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.From, tmpl)

	svc := service.New(mail, rdb)

	// Servidor gRPC.
	grpcServer := grpc.NewServer()
	notificationsv1.RegisterNotificationServiceServer(grpcServer, notifgrpc.NewServer(svc))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Error("no se pudo escuchar", "port", cfg.GRPCPort, "err", err)
		os.Exit(1)
	}

	go func() {
		log.Info("servicio notifications escuchando", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("apagando servicio notifications...")
	grpcServer.GracefulStop()
}
