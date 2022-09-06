package users

import (
	"time"

	"github.com/erviangelar/go-users-api/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type handler struct {
	DB *sqlx.DB
}

type UserResponse struct {
	ID        uint      `example:"1" format:"int64"`
	Username  string    `json:"username" example:"admin"`
	Name      string    `json:"name" example:"Admin"`
	Role      string    `json:"role" example:"admin"`
	CreatedAt time.Time `json:"created_at" example:"04/09/2022"`
}

type UserRequestBody struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type UpdateUserRequestBody struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func RegisterRoutes(r *gin.RouterGroup, db *sqlx.DB) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/users").Use(middleware.AccessToken())
	routes.GET("", h.GetUsers)
	routes.GET("/:id", h.GetUser)
	routesAdmin := r.Group("/user").Use(middleware.AdminValidate())
	routesAdmin.POST("", h.AddUser)
	routesAdmin.PUT("/:id", h.UpdateUser)
	routesAdmin.DELETE("/:id", h.DeleteUser)
}
