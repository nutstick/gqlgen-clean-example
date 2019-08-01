package postgresql

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Target is parameters to get all PostgresSQLConnection's dependencies
type Target struct {
	fx.In
	PostgresQLURL string `name:"postgresql_url" optional:"true"`
	Lc            fx.Lifecycle
	Logger        *zap.Logger
}

// Connection is connection provider to access to global postgres client
type Connection struct {
	client *sql.DB
}

// New is constructor of PostgresSQLConnection
func New(target Target) (*Connection, error) {
	if target.PostgresQLURL == "" {
		return nil, nil
	}

	client, err := sql.Open("postgres", target.PostgresQLURL)
	if err != nil {
		log.Fatal(err)
	}

	target.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			target.Logger.Info("Connecting PostGresQL server at " + target.PostgresQLURL + ".")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			target.Logger.Info("Disconnect PostGresQL server at " + target.PostgresQLURL + ".")
			return client.Close()
		},
	})
	return &Connection{
		client,
	}, err
}

type key string

const (
	// postgresQLClient key for postgres session in each request
	postgresQLClient key = "postgresql_client"
)

// Connect is method return adpater for http request that
// inject the database client in context
func (m *Connection) Connect() gin.HandlerFunc {
	return func(c *gin.Context) {
		// save it in the mux context
		ctx := context.WithValue(c.Request.Context(), postgresQLClient, m.client)
		c.Request = c.Request.WithContext(ctx)
		// pass execution to the original handler
		c.Next()
	}
}

// ForContext is method to get database from context
func ForContext(ctx context.Context) *sql.DB {
	client, ok := ctx.Value(postgresQLClient).(*sql.DB)
	if !ok {
		panic("ctx passing is not contain postgres client")
	}
	return client
}
