package module

import (
	"github.com/nutstick/gqlgen-clean-example/config"
	"github.com/nutstick/gqlgen-clean-example/controller"
	"github.com/nutstick/gqlgen-clean-example/database/mongodb"
	"github.com/nutstick/gqlgen-clean-example/database/postgresql"
	"github.com/nutstick/gqlgen-clean-example/logging"
	"github.com/nutstick/gqlgen-clean-example/server"
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
