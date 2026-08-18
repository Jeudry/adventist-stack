package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/services/notifications/internal/service"
)

type Handler struct {
	svc *service.NotificationService
}

func NewHandler(svc *service.NotificationService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "sendEmail",
		Method:        "POST",
		Path:          "/api/v1/notifications/emails",
		Summary:       "Enviar un correo desde una plantilla",
		Tags:          []string{"notifications"},
		DefaultStatus: 202,
	}, h.SendEmail)

	huma.Register(api, huma.Operation{
		OperationID:   "publishNotification",
		Method:        "POST",
		Path:          "/api/v1/notifications",
		Summary:       "Publicar una notificación",
		Tags:          []string{"notifications"},
		DefaultStatus: 201,
	}, h.Publish)
}

type sendEmailOutput struct{}

func (h *Handler) SendEmail(ctx context.Context, in *httpx.In[SendEmailRequest]) (*sendEmailOutput, error) {
	if err := h.svc.SendEmail(ctx, in.Body.To, in.Body.Template, in.Body.Variables); err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &sendEmailOutput{}, nil
}

func (h *Handler) Publish(ctx context.Context, in *httpx.In[PublishRequest]) (*httpx.Out[PublishResponse], error) {
	id, err := h.svc.Publish(ctx, in.Body.UserID, in.Body.Title, in.Body.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &httpx.Out[PublishResponse]{Body: PublishResponse{NotificationID: id}}, nil
}
