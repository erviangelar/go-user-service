package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Delete User godoc
// @Summary Delete User.
// @Tags Users
// @Accept */*
// @Produce json
// @Param id path int true "User ID"
// @Success 200
// @Router /api/users/{id} [delete]
// @Security ApiKeyAuth
func (h handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// var user models.User
	result, err := h.DB.ExecContext(c, qDelete, id, time.Now())
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	result.LastInsertId()
	// h.DB.Delete(&user)
	c.Status(http.StatusOK)
}

// var qDelete = `DELETE FROM users WHERE id = $1;`
var qDelete = `UPDATE users SET deleted_at=$2 
WHERE id = $1;`
