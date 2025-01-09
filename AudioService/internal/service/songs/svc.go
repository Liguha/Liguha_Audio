package songs

import (
	"AudioService/internal/models"
	"AudioService/internal/models/api"
	"AudioService/internal/models/db"
	"AudioService/internal/ports/events/rpc"
	"context"
	"fmt"
	"strconv"
)

type songsRepo interface {
	AddSong(ctx context.Context, song db.Song) (uint64, error)
}

type Service struct {
	songsRepo     songsRepo
	audioPreparer rpc.AudioPreparerClient
}

func (svc *Service) AddSong(ctx context.Context, song api.CreateSongRequest) error {
	userInfo := models.GetUserFromContext(ctx)
	songAdd := db.Song{
		Name:   song.Name,
		Author: song.Author,
		UserID: strconv.Itoa(int(userInfo.ID)),
		Music:  song.Music,
	}
	songID, err := svc.songsRepo.AddSong(ctx, songAdd)
	if err != nil {
		return err
	}

	request := &rpc.Audio{
		SampleRate: fmt.Sprintf("Song_%d", songID),
	}

	if _, err = svc.audioPreparer.AddAudio(ctx, request); err != nil {
		return err
	}

	return nil
}

func New(songsRepo songsRepo, audioPreparer rpc.AudioPreparerClient) *Service {
	return &Service{songsRepo: songsRepo, audioPreparer: audioPreparer}
}
