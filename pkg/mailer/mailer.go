package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
)

type Mailer struct {
	host      string
	port      int
	user      string
	pass      string
	from      string
	templates *template.Template
}

func New(host string, port int, user, pass, from string, templates *template.Template) *Mailer {
	return &Mailer{
		host:      host,
		port:      port,
		user:      user,
		pass:      pass,
		from:      from,
		templates: templates,
	}
}

func (m *Mailer) Send(to, name string, data map[string]string) error {
	var body bytes.Buffer
	if err := m.templates.ExecuteTemplate(&body, name+".html", data); err != nil {
		return fmt.Errorf("mailer: render %q: %w", name, err)
	}

	subject := data["subject"]
	if subject == "" {
		subject = "Notification"
	}

	msg := m.buildMessage(to, subject, body.String())
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	var auth smtp.Auth
	if m.user != "" {
		auth = smtp.PlainAuth("", m.user, m.pass, m.host)
	}

	if err := smtp.SendMail(addr, auth, m.from, []string{to}, msg); err != nil {
		return fmt.Errorf("mailer: send: %w", err)
	}
	return nil
}

func (m *Mailer) buildMessage(to, subject, htmlBody string) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", m.from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	b.WriteString("\r\n")
	b.WriteString(htmlBody)
	return []byte(b.String())
}
