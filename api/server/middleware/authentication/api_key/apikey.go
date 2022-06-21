package apikey_auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tolopsy/foodpro/api/persistence"
)

type APIKeyAuth struct {
	apiKey    string
	headerKey string
	verifyUser persistence.UserVerifier
}

func NewAPIKeyAuth(key string, verifyUser persistence.UserVerifier) *APIKeyAuth {
	headerKey := "X-API-KEY"
	return &APIKeyAuth{
		apiKey:    key,
		headerKey: headerKey,
		verifyUser: verifyUser,
	}
}

func (auth *APIKeyAuth) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader(auth.headerKey) != auth.apiKey {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong API key provided"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (auth *APIKeyAuth) SignIn(ctx *gin.Context) {
	var user persistence.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while signing in -> " + err.Error()})
		return
	}

	if !auth.verifyUser(user) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or Password"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{auth.headerKey: auth.apiKey})
}

func (auth *APIKeyAuth) SignOut(ctx *gin.Context) {}
