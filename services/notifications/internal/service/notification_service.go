// Package service contiene la lógica del servicio de notificaciones: envío de
// correos con plantillas y publicación de notificaciones en tiempo real por
// Redis pub/sub.
package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/mailer"
	"github.com/Jeudry/adventist-stack/pkg/redis"
)

// channelPrefix es el prefijo de los canales pub/sub por usuario.
const channelPrefix = "notifications:"

// NotificationService envía correos y publica notificaciones.
type NotificationService struct {
	mailer *mailer.Mailer
	redis  *redis.Client
}

// New crea el servicio.
func New(m *mailer.Mailer, r *redis.Client) *NotificationService {
	return &NotificationService{mailer: m, redis: r}
}

// SendEmail renderiza la plantilla indicada y envía el correo.
func (s *NotificationService) SendEmail(_ context.Context, to, template string, vars map[string]string) error {
	return s.mailer.Send(to, template, vars)
}

// Notification es el payload publicado en Redis.
type Notification struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Publish emite una notificación al canal del usuario y devuelve su ID.
func (s *NotificationService) Publish(ctx context.Context, userID, title, body string) (string, error) {
	n := Notification{
		ID:     uuid.NewString(),
		UserID: userID,
		Title:  title,
		Body:   body,
	}

	payload, err := json.Marshal(n)
	if err != nil {
		return "", fmt.Errorf("service: marshal notification: %w", err)
	}

	channel := channelPrefix + userID
	if err := s.redis.Publish(ctx, channel, payload).Err(); err != nil {
		return "", fmt.Errorf("service: publish notification: %w", err)
	}
	return n.ID, nil
}
