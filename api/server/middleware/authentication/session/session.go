package session_auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	redisSessions "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/tolopsy/foodpro/api/persistence"
)

type SessionAuth struct {
	SessionName     string
	Store           sessions.Store
	verifyUser      persistence.UserVerifier
	userIdentifier  string
	sessionTokenKey string
}

func NewSessionAuth(key, address, password string, verifyUser persistence.UserVerifier) (*SessionAuth, error) {
	sessionStore, err := redisSessions.NewStore(10, "tcp", address, password, []byte(key))
	if err != nil {
		return nil, err
	}
	userIdentifier := "username"
	sessionTokenKey := "token"
	sessionName := "user_sessions"
	return &SessionAuth{
		SessionName:     sessionName,
		Store:           sessionStore,
		verifyUser:      verifyUser,
		userIdentifier:  userIdentifier,
		sessionTokenKey: sessionTokenKey,
	}, nil
}

func (sessionAuth *SessionAuth) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionToken := session.Get(sessionAuth.sessionTokenKey)

		if sessionToken == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (sessionAuth *SessionAuth) SignIn(ctx *gin.Context) {
	var user persistence.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while signing in -> " + err.Error()})
		return
	}

	if !sessionAuth.verifyUser(user) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or Password"})
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(ctx)
	session.Set(sessionAuth.userIdentifier, user.Username)
	session.Set(sessionAuth.sessionTokenKey, sessionToken)
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{"message": "User signed in"})
}

func (sessionAuth *SessionAuth) SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{"message": "User signed out"})
}
