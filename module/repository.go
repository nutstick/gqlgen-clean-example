package module

import (
	admin "github.com/nutstick/nithi-backend/packages/admin/repository"
	"go.uber.org/fx"
)

// RepositoriyModule is Repositories fx module
var RepositoriyModule = fx.Provide(
	admin.NewMongoRepository,
)
