package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIKeyAuth struct {
	apiKey string
}

func NewAPIKeyAuth(key string) *APIKeyAuth {
	return &APIKeyAuth{apiKey: key}
}

func (auth *APIKeyAuth) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("X-API-KEY") != auth.apiKey {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong API key provided"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
