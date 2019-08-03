package admin

import (
	"context"

	"github.com/nutstick/nithi-backend/model"
)

// Repository is interface for AdminRepository
type Repository interface {
	GetAll(ctx context.Context) ([]*model.Admin, error)
	GetByID(ctx context.Context, id model.ID) (*model.Admin, error)
	GetByEmail(ctx context.Context, email string) (*model.Admin, error)
	Create(ctx context.Context, admin *model.Admin) (*model.Admin, error)
	Update(ctx context.Context, id model.ID, update *model.Admin) (*model.Admin, error)
	Delete(ctx context.Context) error
	DeleteByID(ctx context.Context, id model.ID) error
}
