package s3client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var s3Client *s3.Client
var s3Bucket string

func InitializeFromEnvironment(ctx context.Context) {
	// TODO: make region configurable!
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	s3Bucket = os.Getenv("S3_BUCKET")
	slog.Info("Connection to S3", "bucket", s3Bucket)

	addr, err := net.LookupHost("s3.eu-north-1.amazonaws.com")
	slog.Info("After LookupHost", "addr", addr)
	if err != nil {
		slog.Info("LookupHost failed", "error", err)
	}
}

func PutString(ctx context.Context, path string, content string) error {
	putInput := &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(path),
		Body:   strings.NewReader(content),
	}

	// Upload to S3
	_, err := s3Client.PutObject(ctx, putInput)
	if err != nil {
		return err
	}
	return nil
}

func GetString(ctx context.Context, path string) (string, error) {
	output, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if ok := errors.As(err, &noSuchKey); ok {
			return "", fmt.Errorf("object %q not found in bucket %q", path, s3Bucket)
		}
		return "", fmt.Errorf("failed to get object: %w", err)
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read object body: %w", err)
	}

	return string(data), nil
}

func Delete(ctx context.Context, path string) error {
	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(path),
	})
	return err
}

func List(ctx context.Context, prefix string, extension string) ([]string, error) {
	var result []string
	slog.Info("s3client.List 1")
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s3Bucket),
		Prefix: aws.String(prefix),
	}

	slog.Info("s3client.List 2")
	paginator := s3.NewListObjectsV2Paginator(s3Client, input)
	slog.Info("s3client.List 3")

	for paginator.HasMorePages() {
		slog.Info("s3client.List 4")
		page, err := paginator.NextPage(ctx)
		slog.Info("s3client.List 5")
		if err != nil {
			slog.Info("s3client.List 5.1", "error", err)
			return nil, fmt.Errorf("listing S3 objects: s3://%s/%s/*.%s %w", s3Bucket, prefix, extension, err)
		}
		slog.Info("s3client.List 6")

		for _, obj := range page.Contents {
			slog.Info("s3client.List 7")
			if strings.HasSuffix(*obj.Key, "."+extension) {
				slog.Info("s3client.List 8")
				result = append(result, *obj.Key)
			}
		}
	}
	slog.Info("s3client.List 9", "result", result)

	return result, nil
}
