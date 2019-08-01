package module

import (
	"github.com/nutstick/nithi-backend/config"
	"github.com/nutstick/nithi-backend/controller"
	"github.com/nutstick/nithi-backend/database/mongodb"
	"github.com/nutstick/nithi-backend/database/postgresql"
	"github.com/nutstick/nithi-backend/logging"
	"github.com/nutstick/nithi-backend/server"
	"go.uber.org/fx"
)

// Module is registry for all module using in application
// will process by fx
var Module = fx.Options(
	fx.Provide(
		config.New,
		logging.New,
		server.New,
		// Database
		mongodb.New,
		postgresql.New,
		// Controller
		controller.NewGraphQLController,
		controller.NewAuth,
	),
	ServiceModule,
	RepositoriyModule,
)
