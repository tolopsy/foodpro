// Handlers unique only to jwt. This handlers are not
// necessarily part of the auth middleware methods.
package jwt_auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
)

func (jwtAuth *JWTAuth) Refresh(ctx *gin.Context) {
	tokenValue := ctx.GetHeader(jwtAuth.headerKey)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, jwtAuth.getTokenSecret)
	
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error while parsing token string -> " + err.Error()})
		return
	}

	if token == nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	expireBoundary := 30 * time.Second
	expiryTime := claims.ExpiresAt.Time
	if time.Until(expiryTime) > expireBoundary {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Token has not expired"})
		return
	}

	expiryTime = time.Now().Add(5 * time.Second)
	claims.ExpiresAt = jwt.NewNumericDate(expiryTime)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err = token.SignedString(jwtAuth.jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtOutput := JWTOutput{
		Token:   tokenValue,
		Expires: expiryTime,
	}

	ctx.JSON(http.StatusOK, jwtOutput)
}
