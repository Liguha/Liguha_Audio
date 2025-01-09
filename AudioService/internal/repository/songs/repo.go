package songs

import (
	"AudioService/internal/models/db"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	master *pgxpool.Pool
	s3     *session.Session

	bucket string
}

func (r *Repo) AddSong(ctx context.Context, song db.Song) (uint64, error) {
	var songID uint64

	query := `
		INSERT INTO music.songs (name, compositor, author_id)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.master.QueryRow(ctx, query, song.Name, song.Author, song.UserID).Scan(&songID)
	if err != nil {
		return 0, err
	}

	uploader := s3manager.NewUploader(r.s3)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fmt.Sprintf("Song_%d", songID)),
		Body:   song.Music,
	})
	if err != nil {
		return 0, err
	}

	return songID, nil
}
func New(master *pgxpool.Pool, s3Repo *session.Session, bucket string) *Repo {
	return &Repo{
		master: master,
		s3:     s3Repo,
		bucket: bucket,
	}
}
