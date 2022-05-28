package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt_auth "github.com/tolopsy/foodpro/server/middleware/authentication/jwt"
	session_auth "github.com/tolopsy/foodpro/server/middleware/authentication/session"
)

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	SignIn(*gin.Context)
	SignOut(*gin.Context)
}

func LoadSpecialFeatures(auth AuthMiddleware, engine *gin.Engine) {
	switch authType := auth.(type) {
	case *jwt_auth.JWTAuth:
		engine.POST("/refresh", authType.Refresh)
	case *session_auth.SessionAuth:
		engine.Use(sessions.Sessions(authType.SessionName, authType.Store))
	}
}
