package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notificationsv1 "github.com/Jeudry/adventist-stack/gen/notifications/v1"
	"github.com/Jeudry/adventist-stack/services/notifications/internal/service"
)

type Server struct {
	notificationsv1.UnimplementedNotificationServiceServer
	svc *service.NotificationService
}

func NewServer(svc *service.NotificationService) *Server {
	return &Server{svc: svc}
}

func (s *Server) SendEmail(ctx context.Context, req *notificationsv1.SendEmailRequest) (*notificationsv1.SendEmailResponse, error) {
	if req.GetTo() == "" || req.GetTemplate() == "" {
		return nil, status.Error(codes.InvalidArgument, "to and template are required")
	}
	if err := s.svc.SendEmail(ctx, req.GetTo(), req.GetTemplate(), req.GetVariables()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &notificationsv1.SendEmailResponse{Sent: true}, nil
}

func (s *Server) PublishNotification(ctx context.Context, req *notificationsv1.PublishNotificationRequest) (*notificationsv1.PublishNotificationResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	id, err := s.svc.Publish(ctx, req.GetUserId(), req.GetTitle(), req.GetBody())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &notificationsv1.PublishNotificationResponse{Published: true, NotificationId: id}, nil
}
