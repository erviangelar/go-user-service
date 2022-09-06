package auth

import (
	"github.com/erviangelar/go-users-api/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type handler struct {
	DB *sqlx.DB
}
type AuthResponse struct {
	ID       uint   `example:"1" format:"int64"`
	Username string `json:"username" example:"admin"`
	Name     string `json:"name" example:"Admin"`
	Role     string `json:"role" example:"admin"`
	// Token        string `json:"token" example:"jwt-token"`
	// RefreshToken string `json:"refresh-token" example:"jwt-refreshtoken"`
}

type AuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterRoutes(r *gin.RouterGroup, db *sqlx.DB) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/auth")
	routes.POST("", h.Login)
	routesRefresh := r.Group("/refresh-token").Use(middleware.RefreshToken())
	routesRefresh.GET("", h.RefershToken)
}
