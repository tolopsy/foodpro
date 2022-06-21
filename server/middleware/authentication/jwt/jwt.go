package jwt_auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tolopsy/foodpro/persistence"
)

type JWTAuth struct {
	jwtSecret  string
	headerKey  string
	verifyUser persistence.UserVerifier
}

func (jwtAuth *JWTAuth) getTokenSecret(token *jwt.Token) (interface{}, error) {
	return []byte(jwtAuth.jwtSecret), nil
}

func NewJWTAuth(secret string, verifyUser persistence.UserVerifier) *JWTAuth {
	headerKey := "Authorization"
	return &JWTAuth{
		jwtSecret:  secret,
		headerKey:  headerKey,
		verifyUser: verifyUser,
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func (jwtAuth *JWTAuth) SignIn(ctx *gin.Context) {
	var user persistence.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while signing in -> " + err.Error()})
		return
	}

	if !jwtAuth.verifyUser(user) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or Password"})
		return
	}

	expiresAt := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtAuth.jwtSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Error while signing jwt token -> "+err.Error())
		return
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expiresAt,
	}

	ctx.JSON(http.StatusOK, jwtOutput)
}

func (jwtAuth *JWTAuth) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenValue := ctx.GetHeader(jwtAuth.headerKey)
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenValue, claims, jwtAuth.getTokenSecret)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error while parsing token ->" + err.Error()})
			ctx.Abort()
			return
		}

		if token == nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// since I'm not persisting tokens on the server (yet), I'll leave this as dummy.
func (jwtAuth *JWTAuth) SignOut(ctx *gin.Context) {}
