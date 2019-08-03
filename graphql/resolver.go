package graphql

import (
	"context"
	"errors"

	"github.com/nutstick/nithi-backend/constant"
	"github.com/nutstick/nithi-backend/model"
	"github.com/nutstick/nithi-backend/packages/admin"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var (
	// ErrIncorrectPassword throw when mutation Login incorrect password
	ErrIncorrectPassword = errors.New("Incorrect password")
)

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*LoginPayload, error) {
	a, err := admin.ForContext(ctx).GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Verify input password with stored password
	if !admin.ForContext(ctx).ComparePassword(ctx, a, password) {
		return nil, ErrIncorrectPassword
	}

	session := ctx.Value(constant.Session).(*string)
	*session = string(a.ID)

	return &LoginPayload{a}, err
}

func (r *mutationResolver) Register(ctx context.Context, input RegisterInput, secret *string) (*RegisterPayload, error) {
	a, err := admin.ForContext(ctx).Create(ctx, &model.Admin{
		Email:    input.Email,
		Password: input.Password,
		Name:     input.Name,
		Avatar:   input.Avatar,
		Roles:    input.Roles,
	})
	if err != nil {
		return nil, err
	}

	session := ctx.Value(constant.Session).(*string)
	*session = string(a.ID)

	return &RegisterPayload{a}, err
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Helloworld(ctx context.Context) (string, error) {
	return "Helloworld", nil
}
func (r *queryResolver) Viewer(ctx context.Context) (*model.Admin, error) {
	session := ctx.Value(constant.Session).(*model.ID)
	if session == nil {
		return nil, nil
	}
	return admin.ForContext(ctx).GetByID(ctx, *session)
}
