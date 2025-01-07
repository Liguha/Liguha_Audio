package config

import "fmt"

// HTTPServer holds the configuration for the HTTP server.
type HTTPServer struct {
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required"`

	AudioURL string `yaml:"audio_url" validate:"required"`
	AlbumURL string `yaml:"album_url" validate:"required"`
}

func (h HTTPServer) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

type DB struct {
	Address string `yaml:"address" validate:"required"`

	UserName string `yaml:"user_name" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	DataBase string `yaml:"data_base" validate:"required"`
}

type Config struct {
	// HTTPServer is the HTTP server configuration.
	HTTPServer HTTPServer `yaml:"http_server" validate:"required"`

	DB DB `yaml:"db" validate:"required"`

	JWTSecret string `yaml:"jwt_secret"`
}
