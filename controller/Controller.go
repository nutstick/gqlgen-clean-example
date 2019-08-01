package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Controller is interface for Controller
type Controller interface {
	Register(mux *gin.Engine)
}

// Target is parameters to get mux, controllers
type Target struct {
	fx.In
	Gin         *gin.Engine
	Controllers []Controller `group:"controller"`
}

// Result is return results for controllers
type Result struct {
	fx.Out
	Controller Controller `group:"controller"`
}
