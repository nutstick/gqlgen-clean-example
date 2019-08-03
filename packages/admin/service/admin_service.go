package service

import (
	"context"

	"github.com/nutstick/nithi-backend/model"
	"github.com/nutstick/nithi-backend/packages/admin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// adminService handles the creation, modification and deletion of admins.
// It uses a AdminRepository to communicate with the database.
type adminService struct {
	admin.Repository
	logger *zap.Logger
}

// Target is parameters to get all dependencies
type Target struct {
	fx.In
	Repository admin.Repository
	Logger     *zap.Logger
}

// NewService is adminService's constructor
func NewService(target Target) admin.Service {
	return &adminService{
		Repository: target.Repository,
		logger:     target.Logger,
	}
}

// Register new admin and stored it to database with hashed password
func (m *adminService) Register(ctx context.Context, admin *model.Admin) (*model.Admin, error) {
	return m.Repository.Create(ctx, admin)
}

// ComparePassword to compare cryped password
func (m *adminService) ComparePassword(ctx context.Context, admin *model.Admin, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	return err == nil
}
