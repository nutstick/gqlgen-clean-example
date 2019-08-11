package controller

import (
	context "context"
	http "net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nutstick/gqlgen-clean-example/constant"
	"github.com/nutstick/gqlgen-clean-example/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type key string

type authResponseWriter struct {
	gin.ResponseWriter
	tokenPassword     string
	sessionToResolver string
	sessionFromCookie string
}

func (m *authResponseWriter) Write(data []byte) (n int, err error) {
	if m.sessionToResolver != m.sessionFromCookie {
		tk := &model.Token{UserID: m.sessionToResolver}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte(m.tokenPassword))
		http.SetCookie(m, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			HttpOnly: true,
		})
	}
	return m.ResponseWriter.Write(data)
}

// Auth is class provided middleware for read/write token in to request/response cache
type Auth struct {
	tokenPassword string
	logger        *zap.Logger
}

// AuthTarget is parameter object for geting all Auth's dependency
type AuthTarget struct {
	fx.In
	TokenPassword string `name:"token_password"`
	Logger        *zap.Logger
}

// NewAuth is construct for Auth
func NewAuth(target AuthTarget) *Auth {
	return &Auth{
		tokenPassword: target.TokenPassword,
		logger:        target.Logger,
	}
}

// Middleware for GraphQL resolver to pass session into ctx and receive
// session from ctx set in request.Cache
func (m *Auth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		arw := &authResponseWriter{c.Writer, m.tokenPassword, "", ""}

		tokenPart, _ := c.Cookie("token")
		if tokenPart != "" {
			tk := &model.Token{}

			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(m.tokenPassword), nil
			})
			// Malformed token
			if err != nil {
				m.logger.Error("malformed token", zap.Error(err))
			} else {
				// Token is invalid, maybe not signed on this server
				if token.Valid {
					arw.sessionFromCookie = tk.UserID
					arw.sessionToResolver = tk.UserID
				}
			}
		}
		ctx := context.WithValue(c.Request.Context(), constant.Session, &arw.sessionToResolver)
		c.Request = c.Request.WithContext(ctx)
		c.Writer = arw

		c.Next()
	}
}
