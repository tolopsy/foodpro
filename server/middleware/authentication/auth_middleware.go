package auth

import "github.com/gin-gonic/gin"

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	SignIn(*gin.Context)
}