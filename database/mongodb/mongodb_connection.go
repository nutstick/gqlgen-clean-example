package mongodb

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Target is parameters to get all MongoDBConnection's dependencies
type Target struct {
	fx.In
	MongoURL string `name:"mongo_url" optional:"true"`
	Lc       fx.Lifecycle
	Logger   *zap.Logger
}

// Connection is connection provider to access to global mongodb client
type Connection struct {
	client *mongo.Client
}

// New is constructor of MongoDBConnection
func New(target Target) (*Connection, error) {
	if target.MongoURL == "" {
		return nil, nil
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(target.MongoURL))

	target.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			target.Logger.Info("Connecting MongoDB server at " + target.MongoURL + ".")
			return client.Connect(ctx)
		},
		OnStop: func(ctx context.Context) error {
			target.Logger.Info("Disconnect MongoDB server at " + target.MongoURL + ".")
			return client.Disconnect(ctx)
		},
	})
	return &Connection{
		client,
	}, err
}

type key string

const (
	// mongoClient key for mongo session in each request
	mongoClient key = "mongo_client"
)

// Connect is method return adpater for http request that
// inject the database client in context
func (m *Connection) Connect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m != nil {
			// save it in the mux context
			ctx := context.WithValue(c.Request.Context(), mongoClient, m.client)
			c.Request = c.Request.WithContext(ctx)
		} else {
			// TODO: Warn
		}
		// pass execution to the original handler
		c.Next()
	}
}

// WithContext is method apply mongoClient into context
func (m *Connection) WithContext(ctx context.Context) context.Context {
	if m != nil {
		// save it in the mux context
		return context.WithValue(ctx, mongoClient, m.client)
	} else {
		// TODO: Warn
	}
	return ctx
}

// ForContext is method to get mongodb client from context
func ForContext(ctx context.Context) *mongo.Client {
	client, ok := ctx.Value(mongoClient).(*mongo.Client)
	if !ok {
		panic("ctx passing is not contain mongodb client")
	}
	return client
}
