package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	contextUser             = "context_user"
)

func (s *server) applyAuthentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Vary", authorizationHeaderKey)

		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token type",
			})
			return
		}

		token := headerParts[1]
		payload, err := s.jwt.VerifyToken(token)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "token expired",
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		user, err := s.store.FindUserByID(payload.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		ctx.Set(contextUser, *user)
		ctx.Next()
	}
}

func (s *server) getAuthUser(ctx *gin.Context) domain.User {
	user, ok := ctx.Get(contextUser)
	if !ok {
		panic("missing user value in context")
	}
	return user.(domain.User)
}
