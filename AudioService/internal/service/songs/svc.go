package songs

import (
	"AudioService/internal/models"
	"AudioService/internal/models/api"
	"AudioService/internal/models/db"
	"context"
	"strconv"
)

type songsRepo interface {
	AddSong(ctx context.Context, song db.Song) (uint64, error)
}

type Service struct {
	songsRepo songsRepo
}

func (svc *Service) AddSong(ctx context.Context, song api.CreateSongRequest) (uint64, error) {
	userInfo := models.GetUserFromContext(ctx)
	songAdd := db.Song{
		Name:   song.Name,
		Author: song.Author,
		UserID: strconv.Itoa(int(userInfo.ID)),
		Music:  song.Music,
	}
	return svc.songsRepo.AddSong(ctx, songAdd)

}

func New(songsRepo songsRepo) *Service {
	return &Service{songsRepo: songsRepo}
}
