package users

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/erviangelar/go-users-api/pkg/common/helper"
	"github.com/gin-gonic/gin"
)

type PageParam struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// Get Users godoc
// @Summary Get Users.
// @Description Get Users.
// @Tags Users
// @Accept */*
// @Produce json
// @Success 200 {object} []UserResponse
// @Router /api/users [get]
// @Security ApiKeyAuth
func (h handler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	role, _ := ctx.Value("role").(string)
	fmt.Println(role)

	var req PageParam
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	pagination := helper.Pagination(c)
	offset := (pagination.Page - 1) * pagination.Limit
	roleUser := ""
	if strings.ToLower(role) == "user" {
		roleUser = "AND role ='user'"
	}
	fmt.Println(roleUser)
	qSelectList, qCountList := generateQuery(roleUser)
	rows, err := h.DB.QueryContext(c, qSelectList, pagination.Limit, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer rows.Close()
	var users []UserResponse
	for rows.Next() {
		var user UserResponse
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Role,
			&user.CreatedAt,
		); err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
		users = append(users, user)
	}
	if err := rows.Close(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	var total int
	row := h.DB.QueryRowContext(c, qCountList)
	if err := row.Scan(
		&total,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &users, "page": pagination.Page, "limit": pagination.Limit, "total": total})
}

func generateQuery(roleUser string) (string, string) {

	qSelectList := `SELECT "id","name","username", "role","created_at" 
	FROM users WHERE deleted_at is null ` + roleUser + ` 
	LIMIT $1
	OFFSET $2`
	qCountList := `SELECT COUNT("id") as total 
	FROM users WHERE deleted_at is null ` + roleUser
	return qSelectList, qCountList
}
