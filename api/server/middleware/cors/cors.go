package cors_middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewCorsMiddleware(rule cors.Config) gin.HandlerFunc {
	return cors.New(rule)
}