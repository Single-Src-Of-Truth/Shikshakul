package email

import (
	"context"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"go.uber.org/zap"
)

type SESProvider struct {
	client *sesv2.Client
	sender string
	logger *zap.Logger
}

func NewSESProvider(ctx context.Context, region string, senderEmail string, logger *zap.Logger) (*SESProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &SESProvider{
		client: sesv2.NewFromConfig(cfg),
		sender: senderEmail,
		logger: logger,
	}, nil
}

func (p *SESProvider) SendEmail(ctx context.Context, payload domain.EmailPayload) error {
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(p.sender),
		Destination: &types.Destination{
			ToAddresses: payload.To,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(payload.Subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(payload.Body),
					},
				},
			},
		},
	}

	_, err := p.client.SendEmail(ctx, input)
	if err != nil {
		p.logger.Error("Failed to dispatch email via SES",
			zap.Error(err),
			zap.Strings("to", payload.To),
			zap.String("subject", payload.Subject),
		)
		return err
	}

	return nil
}

func (p *SESProvider) Ping(ctx context.Context) error {
	_, err := p.client.GetAccount(ctx, &sesv2.GetAccountInput{})
	return err
}
