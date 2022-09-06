package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/erviangelar/go-users-api/pkg/common/jwt"
	"github.com/gin-gonic/gin"
)

func AccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have an access"})
			ctx.Abort()
			return
		}
		role, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}
		// fmt.Println(role)
		// ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "role", role))
		// c := context.WithValue(ctx.Request.Context(), "role", &role)
		// ctx.Request = ctx.Request.WithContext(c)
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "role", role))
		ctx.Next()
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		// fmt.Println(tokenString)
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "request does not contain an refresh token"})
			ctx.Abort()
			return
		}
		ID, err := jwt.ValidateRefreshToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}
		// fmt.Println(role)
		// ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "role", role))
		// c := context.WithValue(ctx.Request.Context(), "role", &role)
		// ctx.Request = ctx.Request.WithContext(c)
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "ID", ID))
		ctx.Next()
	}
}

func AdminValidate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have an access"})
			ctx.Abort()
			return
		}
		role, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}
		if strings.ToLower(role) != "admin" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "your role don't have an access"})
			ctx.Abort()
			return
		}
		// fmt.Println(role)
		// ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "role", role))
		// c := context.WithValue(ctx.Request.Context(), "role", &role)
		// ctx.Request = ctx.Request.WithContext(c)
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "role", role))
		ctx.Next()
	}
}
