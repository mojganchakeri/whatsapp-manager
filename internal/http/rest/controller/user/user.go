package user

import (
	"github.com/mojganchakeri/whatsapp-manager/internal/init/http"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/logger"
	"github.com/mojganchakeri/whatsapp-manager/internal/service/user"
	"github.com/mojganchakeri/whatsapp-manager/internal/service/whatsapp"
)

type UserController interface {
	Login(c *http.HttpContext) error
	GetStatus(c *http.HttpContext) error
}

type userController struct {
	log             logger.Logger
	whatsappService whatsapp.WhatsappService
	userService     user.UserService
}

func New(
	log logger.Logger,
	whatsappService whatsapp.WhatsappService,
	userService user.UserService,
) UserController {
	return &userController{
		log:             log,
		whatsappService: whatsappService,
		userService:     userService,
	}
}

func (uc *userController) Login(c *http.HttpContext) error {
	userID := c.FormValue("user_id")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(response{Message: "user_id is empty"})
	}

	client, err := uc.whatsappService.CreateClient(c.Context(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	code, err := uc.whatsappService.Login(c.Context(), client, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response{Message: err.Error()})
	} else if code == "" {
		return c.Status(http.StatusOK).JSON(response{Message: "user is already logged in"})
	}

	c.Status(http.StatusOK).JSON(response{
		Message: "scan the QR code",
		QRCode:  code,
	})

	return nil
}

func (uc *userController) GetStatus(c *http.HttpContext) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(response{Message: "user_id is empty"})
	}

	status, err := uc.userService.GetUserStatusByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	c.Status(http.StatusOK).JSON(response{
		Status: status,
	})

	return nil
}
