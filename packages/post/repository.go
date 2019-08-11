package post

import (
	"context"

	"github.com/nutstick/gqlgen-clean-example/model"
)

// Repository is interface for PostRepository
type Repository interface {
	GetAll(ctx context.Context) ([]*model.Post, error)
	GetByID(ctx context.Context, id model.ID) (*model.Post, error)
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	Update(ctx context.Context, id model.ID, update *model.Post) (*model.Post, error)
	Delete(ctx context.Context) error
	DeleteByID(ctx context.Context, id model.ID) error
}
