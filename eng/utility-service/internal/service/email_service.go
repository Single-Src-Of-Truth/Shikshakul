package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/events"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type EmailService struct {
	provider domain.EmailProvider
	limiter  *rate.Limiter
	logger   *zap.Logger
}

func NewEmailService(provider domain.EmailProvider, reqPerSec int, logger *zap.Logger) *EmailService {
	limiter := rate.NewLimiter(rate.Limit(reqPerSec), reqPerSec)

	return &EmailService{
		provider: provider,
		limiter:  limiter,
		logger:   logger,
	}
}

func (s *EmailService) ProcessEmailEvent(ctx context.Context, event events.EmailEvent) error {
	if err := s.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limiter context canceled: %w", err)
	}

	renderedHTML, subject, err := s.renderTemplate(event.Type, event.Data)
	if err != nil {
		s.logger.Error("Failed to render email template", zap.String("type", event.Type), zap.Error(err))
		return err
	}

	payload := domain.EmailPayload{
		To:      []string{event.Recipient},
		Subject: subject,
		Body:    renderedHTML,
	}

	if err := s.provider.SendEmail(ctx, payload); err != nil {
		return fmt.Errorf("provider dispatch failed: %w", err)
	}

	s.logger.Info("Email successfully dispatched", zap.String("event_id", event.EventID), zap.String("recipient", event.Recipient))
	return nil
}

func (s *EmailService) renderTemplate(eventType string, data map[string]interface{}) (html string, subject string, err error) {
	var templatePath string
	var tmplData interface{}

	switch eventType {
	case "OTP_VERIFICATION":
		templatePath = "internal/templates/html/otp.html"
		subject = "Your Shikshakul Security Code"

		otp, _ := data["otp_code"].(string)
		tmplData = struct {
			OTPCode string
			Date    string
		}{
			OTPCode: otp,
			Date:    time.Now().Format("January 2, 2006"),
		}

	case "INVITE":
		templatePath = "internal/templates/html/verification.html"
		subject = "Join Your Shikshakul Workspace"

		link, _ := data["invite_link"].(string)
		tmplData = struct {
			InviteLink string
		}{
			InviteLink: link,
		}

	case "GENERIC":
		templatePath = "internal/templates/html/generic.html"

		subj, _ := data["subject"].(string)
		body, _ := data["message"].(string)
		subject = subj

		tmplData = struct {
			Subject string
			Body    string
		}{
			Subject: subj,
			Body:    body,
		}

	default:
		return "", "", fmt.Errorf("unknown email event type: %s", eventType)
	}

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse template file %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, tmplData); err != nil {
		return "", "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), subject, nil
}
