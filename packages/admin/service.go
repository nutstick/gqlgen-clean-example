package admin

import (
	"context"

	"github.com/nutstick/gqlgen-clean-example/model"
)

type key string

const (
	// Admin key for admin service in each request
	Key key = "admin"
)

// Service is interface for AdminService
type Service interface {
	Repository
	Register(ctx context.Context, admin *model.Admin) (*model.Admin, error)
	ComparePassword(ctx context.Context, admin *model.Admin, password string) bool
}

// ForContext is method to get admin service from context
func ForContext(ctx context.Context) Service {
	service, ok := ctx.Value(Key).(Service)
	if !ok {
		panic("ctx passing is not contain admin service")
	}
	return service
}
