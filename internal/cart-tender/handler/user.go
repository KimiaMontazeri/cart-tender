package handler

import (
	"errors"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/auth"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/model"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/request"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	userRepo      model.UserRepo
	authenticator auth.Auth
}

func NewUserHandler(userRepo model.UserRepo, authenticator auth.Auth) *UserHandler {
	return &UserHandler{userRepo: userRepo, authenticator: authenticator}
}

func (h *UserHandler) Register(c echo.Context) error {
	req := new(request.UserRequest)
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("register: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user := model.NewUser(req.Username, req.Password)

	if err := h.userRepo.Create(user); err != nil {
		if errors.Is(err, model.ErrUserExists) {
			return c.String(http.StatusBadRequest, err.Error())
		}

		logrus.Infof("register: failed to create: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	token, err := h.authenticator.GenerateUserJWT(*user)
	if err != nil {
		logrus.Infof("register: failed to generate jwt: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(request.UserRequest)
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("login: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := h.userRepo.Find(req.Username)
	if err != nil {
		logrus.Infof("login: failed to find: %s", err.Error())
		return c.NoContent(http.StatusForbidden)
	}

	if !user.CheckPassword(req.Password) {
		logrus.Info("login: incorrect password")
		return c.NoContent(http.StatusForbidden)
	}

	token, err := h.authenticator.GenerateUserJWT(user)
	if err != nil {
		logrus.Infof("login: failed to generate jwt: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
