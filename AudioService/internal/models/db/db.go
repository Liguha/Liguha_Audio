package db

import "bytes"

type Song struct {
	Name       string `json:"name"`
	Author     string `json:"author"`
	Tags       string `json:"tags"`
	UserID     string `json:"user_id"`
	IsOfficial bool   `json:"is_official"`
	Music      *bytes.Buffer
}
