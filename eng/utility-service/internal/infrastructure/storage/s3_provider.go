package storage

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"go.uber.org/zap"
)

type S3Provider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
	region        string
	logger        *zap.Logger
}

func NewS3Provider(ctx context.Context, region, bucketName string, logger *zap.Logger) (*S3Provider, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return &S3Provider{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    bucketName,
		region:        region,
		logger:        logger,
	}, nil
}

func (p *S3Provider) InitializeInfrastructure(ctx context.Context) error {
	p.logger.Info("Initializing S3 Infrastructure...")

	_, err := p.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(p.bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(p.region),
		},
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "BucketAlreadyOwnedByYou" || apiErr.ErrorCode() == "BucketAlreadyExists" {
				p.logger.Info("S3 Bucket already exists, proceeding to policy updates.")
			} else {
				return err
			}
		}
	} else {
		p.logger.Info("Successfully created new S3 Bucket", zap.String("bucket", p.bucketName))
	}

	_, err = p.client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
		Bucket: aws.String(p.bucketName),
		CORSConfiguration: &types.CORSConfiguration{
			CORSRules: []types.CORSRule{
				{
					AllowedHeaders: []string{"*"},
					AllowedMethods: []string{"PUT", "POST", "GET"},
					AllowedOrigins: []string{"*"}, // TODO: Restrict in prod
					ExposeHeaders:  []string{"ETag"},
					MaxAgeSeconds:  aws.Int32(3000),
				},
			},
		},
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
			p.logger.Warn("Access Denied when applying CORS. Ensure the bucket has correct CORS if needed.", zap.Error(err))
		} else {
			p.logger.Error("Failed to apply CORS", zap.Error(err))
			return err
		}
	}

	_, err = p.client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(p.bucketName),
		LifecycleConfiguration: &types.BucketLifecycleConfiguration{
			Rules: []types.LifecycleRule{
				{
					ID:     aws.String("Cleanup-Temp-Folder"),
					Status: types.ExpirationStatusEnabled,
					Prefix: aws.String("temp/"),
					Expiration: &types.LifecycleExpiration{
						Days: aws.Int32(7),
					},
				},
				{
					ID:     aws.String("Cleanup-Trash-Folder"),
					Status: types.ExpirationStatusEnabled,
					Prefix: aws.String("trash/"),
					Expiration: &types.LifecycleExpiration{
						Days: aws.Int32(60),
					},
				},
			},
		},
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
			p.logger.Warn("Access Denied when applying Lifecycle Rules. Ensure the bucket has correct rules if needed.", zap.Error(err))
		} else {
			p.logger.Error("Failed to apply Lifecycle Rules", zap.Error(err))
			return err
		}
	}

	p.logger.Info("S3 Infrastructure fully configured.")
	return nil
}

func (p *S3Provider) GenerateUploadURL(ctx context.Context, objectKey string, contentType string) (string, error) {
	req, err := p.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (p *S3Provider) GenerateDownloadURL(ctx context.Context, objectKey string, inline bool) (string, error) {
	disposition := "attachment"
	if inline {
		disposition = "inline"
	}

	req, err := p.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(p.bucketName),
		Key:                        aws.String(objectKey),
		ResponseContentDisposition: aws.String(disposition),
	}, s3.WithPresignExpires(1*time.Hour))

	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (p *S3Provider) MoveObject(ctx context.Context, sourceKey, destKey string) error {
	_, err := p.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(p.bucketName),
		CopySource: aws.String(p.bucketName + "/" + sourceKey),
		Key:        aws.String(destKey),
	})
	if err != nil {
		return err
	}
	return p.DeleteObject(ctx, sourceKey)
}

func (p *S3Provider) DeleteObject(ctx context.Context, objectKey string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(objectKey),
	})
	return err
}

func (p *S3Provider) Ping(ctx context.Context) error {
	_, err := p.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(p.bucketName),
	})
	return err
}
