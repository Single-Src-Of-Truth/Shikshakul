package domain

import "context"

type EmailPayload struct {
	To      []string
	Subject string
	Body    string
}

type EmailProvider interface {
	SendEmail(ctx context.Context, payload EmailPayload) error
	Ping(ctx context.Context) error
}
