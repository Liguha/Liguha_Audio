package config

import "fmt"

// HTTPServer holds the configuration for the HTTP server.
type HTTPServer struct {
	Host string `envconfig:"HOST" validate:"required"`
	Port uint64 `envconfig:"PORT" validate:"required"`
}

func (h HTTPServer) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

type DB struct {
	Host string `envconfig:"DB_HOST" validate:"required"`
	Port uint64 `envconfig:"DB_PORT" validate:"required"`

	UserName string `envconfig:"DB_USER_NAME" validate:"required"`
	Password string `envconfig:"DB_PASSWORD" validate:"required"`
	DataBase string `envconfig:"DB_NAME" validate:"required"`
}

func (d DB) Address() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}

type S3 struct {
	EndPoint   string `envconfig:"S3_ENDPOINT" validate:"required"`
	Region     string `envconfig:"S3_REGION" validate:"required"`
	KeyID      string `envconfig:"S3_KEY_ID" validate:"required"`
	KeySecret  string `envconfig:"S3_KEY_SECRET" validate:"required"`
	BucketName string `envconfig:"S3_BUCKET_NAME" validate:"required"`
}

type Config struct {
	HTTPServer HTTPServer
	DB         DB
	S3         S3
}
