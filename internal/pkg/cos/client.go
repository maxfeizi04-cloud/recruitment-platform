package cos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/config"

	cosapi "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/google/uuid"
)

type Uploader interface {
	Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error)
}

type TencentCOS struct {
	client    *cosapi.Client
	bucketURL string
}

func NewTencentCOS(cfg config.COSConfig) (Uploader, error) {
	u, err := url.Parse(cfg.BucketURL)
	if err != nil {
		return nil, fmt.Errorf("parse bucket url: %w", err)
	}

	b := &cosapi.BaseURL{BucketURL: u}
	client := cosapi.NewClient(b, &http.Client{
		Transport: &cosapi.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	return &TencentCOS{
		client:    client,
		bucketURL: cfg.BucketURL,
	}, nil
}

func (c *TencentCOS) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	opt := &cosapi.ObjectPutOptions{
		ObjectPutHeaderOptions: &cosapi.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	}

	_, err := c.client.Object.Put(ctx, key, reader, opt)
	if err != nil {
		return "", fmt.Errorf("cos put object: %w", err)
	}

	return fmt.Sprintf("%s/%s", c.bucketURL, key), nil
}

func GenerateKey(userID, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return fmt.Sprintf("uploads/%s/%s%s", userID, uuid.New().String()[:8], ext)
}
