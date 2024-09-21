package app

import (
	"fmt"

	userController "github.com/mojganchakeri/whatsapp-manager/internal/http/rest/controller/user"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/config"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/database"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/http"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/logger"
	userConnectionsRepository "github.com/mojganchakeri/whatsapp-manager/internal/repository/postgres/userConnections"
	"github.com/mojganchakeri/whatsapp-manager/internal/router"
	userService "github.com/mojganchakeri/whatsapp-manager/internal/service/user"
	whatsappService "github.com/mojganchakeri/whatsapp-manager/internal/service/whatsapp"

	"go.uber.org/fx"
)

func Run() error {
	app := fx.New(
		fx.Provide(
			// init
			config.New,
			logger.New,
			database.New,
			http.New,

			// postgres repositories
			userConnectionsRepository.New,

			// services
			whatsappService.New,
			userService.New,

			// middlewares

			// controllers
			userController.New,
		),
		fx.Invoke(
			router.Router,
		),
	)

	err := app.Err()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	app.Run()
	return nil
}
