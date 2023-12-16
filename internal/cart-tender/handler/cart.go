package handler

import (
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/auth"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/model"
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/request"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type CartHandler struct {
	cartRepo      model.CartRepo
	authenticator auth.Auth
}

func NewCartHandler(cartRepo model.CartRepo, authenticator auth.Auth) *CartHandler {
	return &CartHandler{cartRepo: cartRepo, authenticator: authenticator}
}

func (h *CartHandler) Create(c echo.Context) error {
	user, ok := c.Get("user").(auth.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	req := new(request.CartCreateRequest)
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("create cart: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	cart := model.NewCart(user.Username, req.Data)
	cart.State = model.PENDING

	if err := h.cartRepo.Create(cart); err != nil {
		logrus.Infof("create cart: failed to create: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *CartHandler) Update(c echo.Context) error {
	user, ok := c.Get("user").(auth.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	req := new(request.CartUpdateRequest)
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("update cart: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	cart := model.Cart{}
	var err error

	if cart, err = h.cartRepo.Find(req.ID); err != nil {
		logrus.Infof("update cart: failed to fetch: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	if user.Username != cart.Username {
		return c.NoContent(http.StatusForbidden)
	}

	cart.Data = req.Data
	cart.State = req.State

	if err := h.cartRepo.Update(&cart); err != nil {
		logrus.Infof("update cart: failed to update: %s", err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CartHandler) Delete(c echo.Context) error {
	user, ok := c.Get("user").(auth.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	req := new(request.CartDeleteRequest)
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("delete cart: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	cart := model.Cart{}
	var err error

	if cart, err = h.cartRepo.Find(req.ID); err != nil {
		logrus.Infof("delete cart: failed to fetch: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	if user.Username != cart.Username {
		return c.NoContent(http.StatusForbidden)
	}

	if err := h.cartRepo.Delete(req.ID); err != nil {
		logrus.Infof("delete cart: failed to delete: %s", err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CartHandler) Find(c echo.Context) error {
	cartID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logrus.Errorf("find cart: %s", err.Error())
		return echo.ErrBadRequest
	}

	user, ok := c.Get("user").(auth.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	cart := model.Cart{}
	if cart, err = h.cartRepo.Find(cartID); err != nil {
		logrus.Infof("find cart: failed to fetch: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	if user.Username != cart.Username {
		return c.NoContent(http.StatusForbidden)
	}

	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) FindByUser(c echo.Context) error {
	user, ok := c.Get("user").(auth.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	carts := []model.Cart{}
	var err error

	if carts, err = h.cartRepo.FindByUser(user.Username); err != nil {
		logrus.Infof("find user carts: failed to fetch: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, carts)
}
