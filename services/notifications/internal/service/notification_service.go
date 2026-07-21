package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/mailer"
	"github.com/Jeudry/adventist-stack/pkg/redis"
)

const channelPrefix = "notifications:"

type NotificationService struct {
	mailer *mailer.Mailer
	redis  *redis.Client
}

func New(m *mailer.Mailer, r *redis.Client) *NotificationService {
	return &NotificationService{mailer: m, redis: r}
}

func (s *NotificationService) SendEmail(_ context.Context, to, template string, vars map[string]string) error {
	return s.mailer.Send(to, template, vars)
}

type Notification struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

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
