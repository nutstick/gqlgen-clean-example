package module

import (
	admin "github.com/nutstick/nithi-backend/packages/admin/service"
	"go.uber.org/fx"
)

// ServiceModule is Repositories fx module
var ServiceModule = fx.Provide(
	admin.NewService,
)
