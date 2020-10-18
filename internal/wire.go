//+build wireinject

package internal

import (
	app "github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/app/svc"
	app_server "github.com/Reasno/kitty/app/svc/server"
	"github.com/google/wire"
)

func InitializeEndpoints() svc.Endpoints {
	wire.Build(app.NewService, app_server.NewEndpoints)
	return svc.Endpoints{}
}




