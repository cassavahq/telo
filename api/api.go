package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"os/signal"
	"syscall"
	"context"
	//Statik
	_ "github.com/cassavahq/telo/statik"
	"github.com/rakyll/statik/fs"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	"github.com/pborman/uuid"
	"github.com/cassavahq/telo/conf"
)

// API is the data holder for the API
type API struct {
	log    *logrus.Entry
	config *conf.Config
	echo   *echo.Echo
	sigFunc     func(os.Signal)
}

//JWTClaims to add custom claims data
type JWTClaims struct {
	jwt.StandardClaims
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Groups []string `json:"groups"`
}

//Valid to check validity token
func (c JWTClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if c.UserID == "" {
		return errors.New("Must provide a user ID")
	}

	return nil
}

func (a *API) waitForSignals() os.Signal {
	sink := make(chan os.Signal, 1)
	defer close(sink)
	// wait for signal
	signal.Notify(sink, signals...)
	// reset the watched signals
	defer signal.Ignore(signals...)

	return <-sink
}

func (a *API) OnSignalReceived(fn func(s os.Signal)) {
	a.sigFunc = fn
}

//Start will start the API on the specified port
func (api *API) Start(){

	go func() {
		api.log.Fatal(api.echo.Start(fmt.Sprintf(":%d", api.config.Port)))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	// main loop
	for {
		// wait for supported signal
		sig := api.waitForSignals()
		api.log.Info("STOPPED")
		// call it once received
		if api.sigFunc != nil {
			api.sigFunc(sig)
		}
		// if it's sighub and rotate the log file when set
		if sig == syscall.SIGHUP {
			api.log.Debug("SIGHUP")
			// TODO: add log rotator
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := api.echo.Shutdown(ctx); err != nil {
			api.log.Fatal(err)
		}
		break
	}
}

// Stop will shutdown the engine internally
func (api *API) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	api.echo.Logger.Fatal(api.echo.Shutdown(ctx))
}

// NewAPI will create an api instance that is ready to start
func NewAPI(log *logrus.Entry, config *conf.Config) *API {
	// create the api
	api := &API{
		config: config,
		log:    log.WithField("component", "api"),
	}

	requireClaims := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: jwt.SigningMethodHS256.Name,
		ContextKey:    tokenKey,
		Claims:        &JWTClaims{},
		SigningKey:    []byte(config.JWTSecret),
	})

	// closed pipeline upon received close signal
	api.OnSignalReceived(func(s os.Signal) {
		api.log.Debug("signal received")
		if s != syscall.SIGHUP {
			// TODO: add close signal 
		}
	})

	// add the endpoints
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	statikFS, err := fs.New()
	if err != nil {
		api.log.Fatal(err)
	}

	staticHandler := http.FileServer(statikFS)

	e.GET("/", echo.WrapHandler(staticHandler))
	e.GET("/statics/*", echo.WrapHandler(staticHandler))

	e.Use(api.setupRequest)
	e.GET("/info", api.Info)
	e.POST("/login", api.generateToken)
	e.GET("/echo", api.dumpToken, requireClaims)

	e.HTTPErrorHandler = api.handleError

	api.echo = e

	return api
}

// Info Handler
func (api *API) Info(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"version":     "testing",
		"description": "a boiler plate project",
		"name":        "telo goreng enakkk",
	})
}

// TokenRequest is the required payload for the generateToken endpoint.
type TokenRequest struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

// TokenResponse defines
type TokenResponse struct {
	Key string `json:"key"`
}

func (api *API) generateToken(ctx echo.Context) error {
	payload := new(TokenRequest)
	if err := ctx.Bind(payload); err != nil {
		return err
	}
	log := getLogger(ctx)

	// validate the payload
	if payload.Email == "" || payload.Pass == "" {
		log.WithFields(logrus.Fields{
			"missing_password": payload.Pass == "",
			"missing_email":    payload.Email == "",
		}).Info("Missing parameters in request")
		return echo.NewHTTPError(http.StatusBadRequest, "Must provide both email and password")
	}
	log.Debug("Starting to issue a new token for a valid request")

	// we have a good payload ~ generate a token
	claims := &JWTClaims{
		UserID: uuid.NewRandom().String(),
		Email:  payload.Email,
	}
	claims.ExpiresAt = time.Now().Add(time.Minute * 60).Unix()

	// create a token with our secret key
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(api.config.JWTSecret))
	if err != nil {
		api.log.WithError(err).Warn("Failed to create a token")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create a token")
	}

	log.Debug("Created a token successfully")
	return ctx.JSON(http.StatusCreated, &TokenResponse{Key: signed})
}

func (api *API) dumpToken(ctx echo.Context) error {
	log := getLogger(ctx)

	token := getToken(ctx)
	claims := token.Claims.(*JWTClaims)

	log.WithFields(logrus.Fields{
		"valid_token":      token.Valid,
		"id":               claims.Id,
		"user_id":          claims.UserID,
		"user_email":       claims.Email,
		"user_groups":      claims.Groups,
		"expires_at":       claims.ExpiresAt,
		"expires_at_human": time.Unix(claims.ExpiresAt, 0).String(),
	}).Info("JWT Token")

	log.Debug("Finished dumping token successfully")
	return nil
}

func (api *API) setupRequest(f echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		logger := api.log.WithFields(logrus.Fields{
			"method":     req.Method,
			"path":       req.URL.Path,
			"request_id": uuid.NewRandom().String(),
		})
		
		ctx.Set(loggerKey, logger)

		startTime := time.Now()
		defer func() {
			rsp := ctx.Response()
			logger.WithFields(logrus.Fields{
				"status_code":  rsp.Status,
				"runtime_nano": time.Since(startTime).Nanoseconds(),
			}).Info("Finished request")
		}()

		logger.WithFields(logrus.Fields{
			"user_agent":     req.UserAgent(),
			"content_length": req.ContentLength,
		}).Info("Starting request")

		// we have to do this b/c if not the final error handler will not
		// in the chain of middleware. It will be called after meaning that the
		// response won't be set properly.
		err := f(ctx)
		if err != nil {
			ctx.Error(err)
		}
		return err
	}
}