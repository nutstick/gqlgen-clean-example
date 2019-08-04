package postgresql

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/fx"
	"go.uber.org/zap"

	// Import Postgress drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	client *gorm.DB
}

// New is constructor of PostgresSQLConnection
func New(target Target) (*Connection, error) {
	if target.PostgresQLURL == "" {
		return nil, nil
	}

	client, err := gorm.Open("postgres", target.PostgresQLURL)
	if err != nil {
		target.Logger.Fatal("postgres connection err: ", zap.Error(err))
	}

	// Zap logger integration
	// client.LogMode(true)
	// client.SetLogger(gormzap.New(target.Logger))

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

func (m *Connection) Client() *gorm.DB {
	return m.client
}

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

// WithContext is method apply database into context
func (m *Connection) WithContext(ctx context.Context) context.Context {
	if m != nil {
		// save it in the mux context
		return context.WithValue(ctx, postgresQLClient, m.client)
	} else {
		// TODO: Warn
	}
	return ctx
}

// ForContext is method to get database from context
func ForContext(ctx context.Context) *gorm.DB {
	client, ok := ctx.Value(postgresQLClient).(*gorm.DB)
	if !ok {
		panic("ctx passing is not contain postgres client")
	}
	return client
}
