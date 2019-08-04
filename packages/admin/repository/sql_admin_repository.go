package repository

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nutstick/nithi-backend/database/postgresql"
	"github.com/nutstick/nithi-backend/model"
	"github.com/nutstick/nithi-backend/packages/admin"
	"github.com/nutstick/nithi-backend/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// sqlRepository contains all the interactions
// with the admin collection stored in sql.
type sqlRepository struct {
	logger *zap.Logger
}

// SQLRepositoryTarget is `fx.In` struct for `fx` to get all dependency to create `AdminSQLRepository`
type SQLRepositoryTarget struct {
	fx.In
	Connection *postgresql.Connection
	Logger     *zap.Logger
}

// NewSQLRepository is AdminRepository's constructor
func NewSQLRepository(target SQLRepositoryTarget) (admin.Repository, error) {
	err := target.Connection.Client().Debug().AutoMigrate(&model.Admin{}).Error
	if err != nil {
		return nil, err
	}
	return &sqlRepository{
		logger: target.Logger,
	}, err
}

// DB method extract database client from context
func (m *sqlRepository) DB(ctx context.Context) *gorm.DB {
	return postgresql.ForContext(ctx)
}

// GetAll returns all the admins stored in the database.
func (m *sqlRepository) GetAll(ctx context.Context) ([]*model.Admin, error) {
	var admins []*model.Admin
	err := m.DB(ctx).Find(&admins).Error
	return admins, err
}

// GetByID returns one admin which is matched by input ID from the database.
func (m *sqlRepository) GetByID(ctx context.Context, id model.ID) (*model.Admin, error) {
	var admin model.Admin
	err := m.DB(ctx).First(&admin, 1).Error
	return &admin, err
}

// GetByEmail returns one admin which is matched by email
func (m *sqlRepository) GetByEmail(ctx context.Context, email string) (*model.Admin, error) {
	var admin model.Admin
	err := m.DB(ctx).Where("email = ?", email).First(&admin).Error
	return &admin, err
}

// Create will insert new admin into database
func (m *sqlRepository) Create(ctx context.Context, admin *model.Admin) (*model.Admin, error) {
	hashedPassword, err := hashPassword(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.Password = hashedPassword
	admin.Roles = []string{}
	admin.CreateAt = time.Now()
	admin.UpdateAt = time.Now()
	err = m.DB(ctx).Debug().Create(admin).Error
	return admin, err
}

// Update will update admin by id
func (m *sqlRepository) Update(ctx context.Context, id model.ID, update *model.Admin) (*model.Admin, error) {
	admin, err := m.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	admin.UpdateAt = time.Now()
	if err := utils.Merge(&admin, *update); err != nil {
		return nil, err
	}
	err = m.DB(ctx).Save(admin).Error
	return admin, err
}

// Delete will remove all admins
func (m *sqlRepository) Delete(ctx context.Context) error {
	return m.DB(ctx).Delete(&model.Admin{}).Error
}

// DeleteByID will remove admin by id from database
func (m *sqlRepository) DeleteByID(ctx context.Context, id model.ID) error {
	return m.DB(ctx).Delete(ctx, model.Admin{ID: id}).Error
}
