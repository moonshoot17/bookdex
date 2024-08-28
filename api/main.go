package main

import (
	"github.com/moonshoot17/bookdex/api/config"
	"github.com/moonshoot17/bookdex/api/database"
	"github.com/moonshoot17/bookdex/api/handlers"
	"github.com/moonshoot17/bookdex/api/middleware"
	"github.com/moonshoot17/bookdex/api/server"
	"github.com/moonshoot17/bookdex/api/storage"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.FxModule,
		database.FxModule,
		storage.FxModule,
		handlers.FxModule,
		middleware.FxModule,
		server.FxModule,
	).Run()
}
