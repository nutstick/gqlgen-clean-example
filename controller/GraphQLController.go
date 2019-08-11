package controller

import (
	context "context"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/nutstick/gqlgen-clean-example/database/mongodb"
	"github.com/nutstick/gqlgen-clean-example/graphql"
	"github.com/nutstick/gqlgen-clean-example/packages/admin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// GraphQLController handle the graphql request, parse request to schema and return results
type GraphQLController struct {
	graphiQLEnable bool
	auth           *Auth
	mongodb        *mongodb.Connection

	admin  admin.Service
	logger *zap.Logger
}

// GraphQLControllerTarget is parameter object for geting all GraphQLController's dependency
type GraphQLControllerTarget struct {
	fx.In
	GraphiQLEnable bool `name:"graphiql_enable"`
	Auth           *Auth
	MongoDB        *mongodb.Connection
	Admin          admin.Service
	Logger         *zap.Logger
}

// NewGraphQLController is a constructor for GraphQLController
func NewGraphQLController(target GraphQLControllerTarget) Result {
	return Result{
		Controller: &GraphQLController{
			graphiQLEnable: target.GraphiQLEnable,
			auth:           target.Auth,
			mongodb:        target.MongoDB,
			admin:          target.Admin,
			logger:         target.Logger,
		},
	}
}

// GrqphQL is defining as the GraphQL handler
func (m *GraphQLController) GrqphQL() gin.HandlerFunc {
	h := handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// GraphiQL is defining as the GraphiQL Page handler
func (m *GraphQLController) GraphiQL() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Register is function to register all controller's endpoint handler
func (m *GraphQLController) Register(r *gin.Engine) {
	r.Use(m.mongodb.Connect()).
		Use(m.Middleware()).
		Use(m.auth.Middleware()).
		POST("/v1/graphql", m.GrqphQL())
	if !m.graphiQLEnable {
		r.GET("/v1/graphiql", m.GraphiQL())
	}
}

// Middleware for GraphQL resolver to pass services into ctx
func (m *GraphQLController) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, admin.Key, m.admin)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
