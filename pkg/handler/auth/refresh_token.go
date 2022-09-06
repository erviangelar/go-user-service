package auth

import (
	"net/http"

	"github.com/erviangelar/go-users-api/pkg/common/jwt"
	"github.com/erviangelar/go-users-api/pkg/common/models"
	"github.com/gin-gonic/gin"
)

// Refresh Token godoc
// @Summary Refresh Token.
// @Description Refresh Token.
// @Tags Auth
// @Accept */*
// @Produce json
// @Success 200 {object} AuthResponse
// @Router /api/refresh-token [get]
// @Security ApiKeyAuth
func (h handler) RefershToken(c *gin.Context) {
	ctx := c.Request.Context()
	ID, _ := ctx.Value("ID").(string)

	var user models.User
	if err := h.DB.GetContext(c, &user, qSelectUser, ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	accessToken, err := jwt.GenerateAccessToken(&user)
	if err != nil {
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
	userRespon := AuthResponse{user.ID, user.Name, user.Role, user.Username}
	c.JSON(http.StatusOK, gin.H{"data": &userRespon, "token": accessToken, "refresh-token": refreshToken})
}

var qSelectUser = `SELECT "id","name","username", "role" FROM users WHERE "id" = $1 AND deleted_at is null;`
