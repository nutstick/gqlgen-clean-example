package service

import (
	"github.com/nutstick/gqlgen-clean-example/packages/post"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// postService handles the creation, modification and deletion of posts.
// It uses a postRepository to communicate with the database.
type postService struct {
	post.Repository
	logger *zap.Logger
}

// Target is parameters to get all dependencies
type Target struct {
	fx.In
	Repository post.Repository
	Logger     *zap.Logger
}

// NewService is postService's constructor
func NewService(target Target) post.Service {
	return &postService{
		Repository: target.Repository,
		logger:     target.Logger,
	}
}
