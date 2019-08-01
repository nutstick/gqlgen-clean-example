package logging

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Target is parameters to get all mux's dependencies
type Target struct {
	fx.In
	Environment string `name:"env"`
}

// New logger used for the whole application.
func New(target Target) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()

	if target.Environment == "live" {
		os.Mkdir("log", os.ModePerm)
		year, month, day := time.Now().Date()
		path := "log/" + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + ".json"
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{path, "stderr"}
	}

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config.Build()
}
