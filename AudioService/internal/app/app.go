package app

import (
	"AudioService/internal/ports/http/songs"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"

	"AudioService/internal/config"
	"AudioService/internal/logger"
	"AudioService/internal/ports/http/middleware"
	songsRepository "AudioService/internal/repository/songs"
	songsService "AudioService/internal/service/songs"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) App {
	return App{cfg: cfg}
}

func (app *App) Run() {
	ctx, cancelProcesses := context.WithCancel(context.Background())
	defer cancelProcesses()

	logger.Init()

	log.
		WithField("host", app.cfg.HTTPServer.Host).
		WithField("port", app.cfg.HTTPServer.Port).
		Info("starting")

	awsSession := app.createSession()

	db := app.initDB(ctx)
	defer db.Close()

	songsRepo := songsRepository.New(db, awsSession, app.cfg.S3.BucketName)
	songsSvc := songsService.New(songsRepo)

	handler := app.InitHandler(songsSvc)

	app.runServer(handler)
}

func (app *App) InitHandler(
	songsSvc *songsService.Service,
) *gin.Engine {
	r := gin.Default()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.WithFields(
			log.Fields{
				"method":       httpMethod,
				"absolutePath": absolutePath,
				"handlerName":  handlerName,
				"nuHandlers":   nuHandlers,
			},
		).Debug()
	}

	r.Use(gin.RecoveryWithWriter(log.StandardLogger().Writer()))
	apiV1 := r.Group("/api/v1", middleware.SimpleMiddleware)

	apiV1.POST("/song", songs.AddSong(songsSvc))

	return r
}

func (app *App) initDB(ctx context.Context) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		app.cfg.DB.UserName, app.cfg.DB.Password, app.cfg.DB.Address(), app.cfg.DB.DataBase)

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	return dbpool
}

func (app *App) runServer(handler http.Handler) {
	const (
		timeout = 5 * time.Second
	)

	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Addr:         app.cfg.HTTPServer.Address(),
	}

	notify := make(chan error, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		notify <- httpServer.ListenAndServe()
		close(notify)
	}()

	select {
	case s := <-interrupt:
		log.Info("get Signal: " + s.String())
	case err := <-notify:
		log.Error(fmt.Errorf("app Run Notify %w", err))
	}

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error(fmt.Errorf("shutdown: %w", err))
	}

	log.Info("server shutdown completed")
}

func (app *App) createSession() *session.Session {
	sess := session.Must(session.NewSession(
		&aws.Config{
			Endpoint: aws.String(app.cfg.S3.EndPoint),
			Region:   aws.String(app.cfg.S3.Region),
			Credentials: credentials.NewStaticCredentials(
				app.cfg.S3.KeyID,
				app.cfg.S3.KeySecret,
				"",
			),
			S3ForcePathStyle: aws.Bool(true), // Используйте, если используете совместимый с S3 сервис
		},
	))
	return sess
}
