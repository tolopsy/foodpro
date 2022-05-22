package auth

import (
	"github.com/gin-gonic/gin"
	jwt_auth "github.com/tolopsy/foodpro/server/middleware/authentication/jwt"
)

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	SignIn(*gin.Context)
}

func LoadSpecialHandlers(auth AuthMiddleware, engine *gin.Engine) {
	switch authType := auth.(type) {
	case *jwt_auth.JWTAuth:
		engine.POST("/refresh", authType.Refresh)
	}
}
