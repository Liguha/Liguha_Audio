package app

import (
	"Gateway/internal/config"
	"Gateway/internal/logger"
	"Gateway/internal/ports/http/middleware"
	"Gateway/internal/ports/http/user"
	userRepo "Gateway/internal/repository/user"
	userService "Gateway/internal/service/user"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

	db := app.initDB(ctx)
	defer db.Close()

	usersRepo := userRepo.New(db)
	userSvc := userService.New(usersRepo, app.cfg.JWTSecret)

	handler := app.InitHandler(userSvc)

	app.runServer(handler)
}

func (app *App) InitHandler(
	userSvc *userService.Svc,
) *gin.Engine {
	r := gin.Default()
	jwtAuth := middleware.JWTAuth(app.cfg.JWTSecret)

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

	r.POST("/register", user.RegisterUser(userSvc))
	r.GET("/login", user.LoginUser(userSvc))

	r.Any("/audio/*path", jwtAuth, createReverseProxy(app.cfg.HTTPServer.AudioURL, "/audio"))
	r.Any("/album/*path", jwtAuth, createReverseProxy(app.cfg.HTTPServer.AlbumURL, "/album"))

	return r
}

func (app *App) initDB(ctx context.Context) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		app.cfg.DB.UserName, app.cfg.DB.Password, app.cfg.DB.Address, app.cfg.DB.DataBase)

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	return dbpool
}

func createReverseProxy(targetURL string, prefix string) func(*gin.Context) {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	return func(c *gin.Context) {
		log.Printf("Request received for %s", c.Request.URL.Path)
		// Modify the request path to exclude the microservice prefix
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, prefix)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
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
