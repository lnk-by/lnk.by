package landingpage

import (
	"context"
	"encoding/json"

	"github.com/lnk.by/aws/s3client"
)

func ReadConfiguration(ctx context.Context, path string, v any) error {
	content, err := s3client.GetString(ctx, path)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(content), v)
	if err != nil {
		return err
	}
	return nil
}

func WriteConfiguration(ctx context.Context, path string, v any) error {
	bytes, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	err = s3client.PutString(ctx, path, string(bytes))
	return err
}
