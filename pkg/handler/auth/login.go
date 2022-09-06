package auth

import (
	"net/http"

	"github.com/erviangelar/go-users-api/pkg/common/jwt"
	"github.com/erviangelar/go-users-api/pkg/common/models"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary Login.
// @Description Login.
// @Tags Auth
// @Accept */*
// @Produce json
// @Param request body AuthRequestBody true "Body"
// @Success 200 {object} AuthResponse
// @Router /api/auth [post]
func (h handler) Login(c *gin.Context) {
	body := AuthRequestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var user models.User
	if err := h.DB.GetContext(c, &user, qSelect, body.Username); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	credErr := user.CheckPassword(body.Password)
	if credErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		c.AbortWithError(http.StatusNotFound, credErr)
		return
	}
	accessToken, err := jwt.GenerateAccessToken(&user)
	if err != nil {
		// data.ToJSON(&GenericError{Error: err.Error()}, w)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to login the user. Please try again later"})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	refreshToken, err := jwt.GenerateRefreshToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to login the user. Please try again later"})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	// tokenString, err := jwt.GenerateJWT(user.Name, user.Username)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	c.Abort()
	// 	return
	// }
	var userRespon AuthResponse
	userRespon.ID = user.ID
	userRespon.Name = user.Name
	userRespon.Role = user.Role
	userRespon.Username = user.Username
	c.JSON(http.StatusOK, gin.H{"data": &userRespon, "token": accessToken, "refresh-token": refreshToken})
}

var qSelect = `SELECT "id","name","username", "role", "password" FROM users WHERE "username" = $1 AND deleted_at is null;`
