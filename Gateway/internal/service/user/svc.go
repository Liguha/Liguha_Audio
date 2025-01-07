package user

import (
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"Gateway/internal/entity"
	"Gateway/internal/models/api"
	"Gateway/internal/models/db"
)

type userRepo interface {
	CreateUser(ctx context.Context, user db.User) error
	GetUserByLogin(ctx context.Context, login string) (db.User, error)
}

type Svc struct {
	userRepo  userRepo
	jwtSecret string
}

func (svc *Svc) RegisterUser(ctx context.Context, userReq api.CreateUserRequest) error {
	hashedPassword, err := hashPassword(userReq.Password)
	if err != nil {
		return err
	}

	userDB := db.User{
		Login:    userReq.Login,
		Password: hashedPassword,
		Role:     userReq.Role,
	}

	return svc.userRepo.CreateUser(ctx, userDB)
}

func (svc *Svc) Login(ctx context.Context, userReq api.CreateUserRequest) (string, error) {
	user, err := svc.userRepo.GetUserByLogin(ctx, userReq.Login)
	if err != nil {
		return "", err
	}

	if ok := comparePasswords(userReq.Password, user.Password); !ok {
		return "", entity.ErrUnauthorized
	}

	return svc.createJWT(user)
}

func (svc *Svc) createJWT(user db.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(int(user.ID)),
		"role":   user.Role,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(svc.jwtSecret))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func comparePasswords(password string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func New(userRepo userRepo, jwtSecret string) *Svc {
	return &Svc{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}
