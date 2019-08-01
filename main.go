package main

import (
	"github.com/nutstick/nithi-backend/controller"
	"github.com/nutstick/nithi-backend/module"
	"go.uber.org/fx"
)

// Register function register all API controllers to Mux
func Register(target controller.Target) {
	for _, controller := range target.Controllers {
		controller.Register(target.Gin)
	}
}

func main() {
	app := fx.New(
		module.Module,
		fx.Invoke(Register),
	)

	app.Run()
}
