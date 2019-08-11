package module

import (
	admin "github.com/nutstick/gqlgen-clean-example/packages/admin/service"
	"go.uber.org/fx"
)

// ServiceModule is Repositories fx module
var ServiceModule = fx.Provide(
	admin.NewService,
)
