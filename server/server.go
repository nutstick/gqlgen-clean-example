package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

// Target is parameters to get all mux's dependencies
type Target struct {
	fx.In
	Environment string `name:"env"`
	Port        string `name:"port"`
	Lc          fx.Lifecycle
	Logger      *zap.Logger
}

// New is constructor to create Mux server on specific addr and port
func New(target Target) *gin.Engine {
	var man *autocert.Manager
	var server *http.Server
	r := gin.New()

	// zap.Logger integration with gin
	r.Use(ginzap.Ginzap(target.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(target.Logger, true))

	if target.Environment != "local" {
		host := ""
		man = &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("certs"),
		}

		server = &http.Server{
			Addr:    host + ":443",
			Handler: r,
			TLSConfig: &tls.Config{
				GetCertificate: man.GetCertificate,
			},
		}
	} else {
		host := "localhost"
		server = &http.Server{
			Addr:    host + ":" + target.Port,
			Handler: r,
		}
	}

	target.Lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if target.Environment != "local" {
				target.Logger.Info("Starting HTTPS server at " + server.Addr)
				go server.ListenAndServeTLS("", "")
				go http.ListenAndServe(":80", man.HTTPHandler(nil))
			} else {
				target.Logger.Info("Starting HTTP server at " + server.Addr)
				go server.ListenAndServe()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			target.Logger.Info("Stopping HTTPS server.")
			return server.Shutdown(ctx)
		},
	})

	return r
}
