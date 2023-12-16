package api

import (
	"errors"
	"fmt"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/auth"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/db"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/handler"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/model"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KimiaMontazeri/cart-tender/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	app := echo.New()

	database := db.WithRetry(db.Connect, cfg.Postgres)

	authenticator := auth.NewAuthenticator(cfg.JWT)

	userRepo := model.NewSQLUserRepo(database)
	userHandler := handler.NewUserHandler(userRepo, authenticator)

	cartRepo := model.NewSQLCartRepo(database)
	cartHandler := handler.NewCartHandler(cartRepo, authenticator)

	app.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	api := app.Group("/api")

	user := api.Group("/user")
	user.POST("/register", userHandler.Register)
	user.POST("/login", userHandler.Login)

	cart := api.Group("/cart")
	cart.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			user, err := authenticator.ValidateUserToken(token)
			if err != nil {
				return false, err
			}

			c.Set("user", user)

			return true, nil
		},
	}))

	cart.POST("", cartHandler.Create)
	cart.PATCH("", cartHandler.Update)
	cart.DELETE("", cartHandler.Delete)
	cart.GET("/:id", cartHandler.Find)
	cart.GET("", cartHandler.FindByUser)

	if err := app.Start(fmt.Sprintf(":%d", cfg.API.Port)); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("echo initiation failed: %s", err)
	}

	logrus.Println("API has been started :D")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// Register API command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "api",
			Short: "Run API to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
