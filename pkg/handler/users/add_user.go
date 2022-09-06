package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

// Add User godoc
// @Summary Add User.
// @Description Add User.
// @Tags Users
// @Accept */*
// @Produce json
// @Param request body UserRequestBody true "Body"
// @Success 200 {object} UserResponse
// @Router /api/users [post]
// @Security ApiKeyAuth
func (h handler) AddUser(c *gin.Context) {
	body := UserRequestBody{}

	// getting request's body
	// if err := c.BindJSON(&body); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := body.HashPassword(body.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var user UserResponse
	err := h.DB.GetContext(
		c,
		&user,
		qInsert,
		&body.Username,
		null.StringFromPtr(&body.Name),
		null.StringFromPtr(&body.Role),
		null.StringFromPtr(&body.Password),
		time.Now(),
	)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, &user)
}

func (user *UserRequestBody) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

var qInsert = `INSERT INTO users(username, name, role, password, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING "id","name","username", "role","created_at" as CreatedAt;`
