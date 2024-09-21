package router

import (
	"context"

	userController "github.com/mojganchakeri/whatsapp-manager/internal/http/rest/controller/user"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/http"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/logger"

	"go.uber.org/fx"
)

func Router(
	lc fx.Lifecycle,
	log logger.Logger,
	httpServer http.Http,
	userController userController.UserController,
) {
	httpServer.ActiveLogger()
	httpServer.ActiveRecover()

	gV1 := httpServer.SetGroup("v1")
	httpServer.SetRoute(gV1, []http.Route{
		{
			Method:  http.Post,
			Path:    "login",
			Handler: userController.Login,
		},
		{
			Method:  http.Get,
			Path:    "status",
			Handler: userController.GetStatus,
		},
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := httpServer.Listen()
				if err != nil {
					log.Error(err.Error())
				}
			}()
			return nil
		},
	})
}
