package post

import (
	"context"
)

type key string

const (
	// Key for post service in each request
	Key key = "post"
)

// Service is interface for PostService
type Service interface {
	Repository
}

// ForContext is method to get post service from context
func ForContext(ctx context.Context) Service {
	service, ok := ctx.Value(Key).(Service)
	if !ok {
		panic("ctx passing is not contain post service")
	}
	return service
}
