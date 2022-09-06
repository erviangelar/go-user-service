package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Edit User godoc
// @Summary Edit User.
// @Tags Users
// @Accept */*
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequestBody true "Body"
// @Success 200 {object} UserResponse
// @Router /api/users/{id} [put]
// @Security ApiKeyAuth
func (h handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	body := UpdateUserRequestBody{}

	// getting request's body
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := h.DB.ExecContext(c, qUpdate, id, body.Name, body.Username, body.Role, time.Now())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	// fmt.Println(&result)
	affected, err := result.RowsAffected()
	fmt.Println(affected)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": `success updated data uid : ` + id})
}

var qUpdate = `UPDATE users SET name=$2,username=$3,role=$4,updated_at=$5 
WHERE id = $1
RETURNING "id","name","username", "role","created_at" as CreatedAt;`
