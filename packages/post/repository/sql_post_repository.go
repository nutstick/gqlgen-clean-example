package repository

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nutstick/gqlgen-clean-example/database/postgresql"
	"github.com/nutstick/gqlgen-clean-example/model"
	"github.com/nutstick/gqlgen-clean-example/packages/post"
	"github.com/nutstick/gqlgen-clean-example/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// sqlRepository contains all the interactions
// with the post collection stored in sql.
type sqlRepository struct {
	logger *zap.Logger
}

// SQLRepositoryTarget is `fx.In` struct for `fx` to get all dependency to create `PostSQLRepository`
type SQLRepositoryTarget struct {
	fx.In
	Connection *postgresql.Connection
	Logger     *zap.Logger
}

// NewSQLRepository is PostRepository's constructor
func NewSQLRepository(target SQLRepositoryTarget) (post.Repository, error) {
	err := target.Connection.Client().Debug().AutoMigrate(&model.Post{}).Error
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

// GetAll returns all the posts stored in the database.
func (m *sqlRepository) GetAll(ctx context.Context) ([]*model.Post, error) {
	var posts []*model.Post
	err := m.DB(ctx).Find(&posts).Error
	return posts, err
}

// GetByID returns one post which is matched by input ID from the database.
func (m *sqlRepository) GetByID(ctx context.Context, id model.ID) (*model.Post, error) {
	var post model.Post
	err := m.DB(ctx).First(&post, 1).Error
	return &post, err
}

// Create will insert new post into database
func (m *sqlRepository) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	post.CreateAt = time.Now()
	post.UpdateAt = time.Now()
	err := m.DB(ctx).Debug().Create(post).Error
	return post, err
}

// Update will update post by id
func (m *sqlRepository) Update(ctx context.Context, id model.ID, update *model.Post) (*model.Post, error) {
	post, err := m.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	post.UpdateAt = time.Now()
	if err := utils.Merge(&post, *update); err != nil {
		return nil, err
	}
	err = m.DB(ctx).Save(post).Error
	return post, err
}

// Delete will remove all posts
func (m *sqlRepository) Delete(ctx context.Context) error {
	return m.DB(ctx).Delete(&model.Post{}).Error
}

// DeleteByID will remove post by id from database
func (m *sqlRepository) DeleteByID(ctx context.Context, id model.ID) error {
	return m.DB(ctx).Delete(ctx, model.Post{ID: id}).Error
}
