package user

import (
	"Gateway/internal/entity"
	"Gateway/internal/models/db"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	master *pgxpool.Pool
}

func (r *Repo) CreateUser(ctx context.Context, user db.User) error {
	query := `
		INSERT INTO music_users.users (login, password, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (login) DO NOTHING`

	_, err := r.master.Exec(ctx, query, user.Login, user.Password, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetUserByLogin(ctx context.Context, login string) (db.User, error) {
	var user db.User

	query := `
	SELECT
		id,
		login,
		password,
		role
	FROM music_users.users
	WHERE login = $1`

	err := pgxscan.Get(ctx, r.master, &user, query, login)
	if err != nil {
		if pgxscan.NotFound(err) {
			return db.User{}, entity.ErrNotFound
		}
		return db.User{}, err
	}

	return user, nil
}

func New(master *pgxpool.Pool) *Repo {
	return &Repo{
		master: master,
	}
}
