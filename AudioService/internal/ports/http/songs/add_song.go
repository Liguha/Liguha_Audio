package songs

import (
	"bytes"
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"AudioService/internal/models/api"
)

type addSongService interface {
	AddSong(ctx context.Context, song api.CreateSongRequest) (uint64, error)
}

func AddSong(songSvc addSongService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx = c.Request.Context()
		)

		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "failed to get file: %s", err.Error())
			return
		}

		buf := new(bytes.Buffer)
		f, err := file.Open()
		if err != nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		if _, err = io.Copy(buf, f); err != nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		req := api.CreateSongRequest{
			Name:   c.Query("name"),
			Author: c.Query("author"),
			Music:  buf,
		}

		if _, err = songSvc.AddSong(ctx, req); err != nil {
			log.Errorf("failed to add song %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusCreated)
	}
}
