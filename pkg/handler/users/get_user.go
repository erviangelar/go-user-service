package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get User godoc
// @Summary Get User.
// @Description Get User.
// @Tags Users
// @Accept */*
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} []UserResponse
// @Router /api/users/{id} [get]
// @Security ApiKeyAuth
func (h handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user UserResponse
	row := h.DB.QueryRowContext(c, qSelect, id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &user})
}

var qSelect = `SELECT "id","name","username", "role","created_at" as CreatedAt FROM users WHERE "id" = $1 AND deleted_at is null;`
