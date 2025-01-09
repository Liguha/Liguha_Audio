package api

import "bytes"

type CreateSongRequest struct {
	Name   string `json:"name"`
	Author string `json:"author"`
	Music  *bytes.Buffer
}
