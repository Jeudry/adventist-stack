package http

type SendEmailRequest struct {
	To        string            `json:"to" required:"true" format:"email"`
	Template  string            `json:"template" required:"true" minLength:"1"`
	Variables map[string]string `json:"variables,omitempty"`
}

type PublishRequest struct {
	UserID string `json:"userId" required:"true" format:"uuid"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type PublishResponse struct {
	NotificationID string `json:"notificationId"`
}
