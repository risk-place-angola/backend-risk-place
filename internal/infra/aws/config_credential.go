package aws

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func LoadDefaultConfigCredentials(ctx context.Context) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			os.Getenv("ACCESS_KEY_ID"),
			os.Getenv("SECRET_ACCESS_KEY"),
			"",
		))),
	)
	if err != nil {
		slog.Error("unable to load AWS SDK config", "error", err)
		return nil, err
	}
	return &cfg, nil
}
