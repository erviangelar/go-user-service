package main

import (
	"log"
	"net/http"

	_ "github.com/erviangelar/go-users-api/docs/userdoc"
	"github.com/erviangelar/go-users-api/pkg/common/config"
	"github.com/erviangelar/go-users-api/pkg/common/db"
	"github.com/erviangelar/go-users-api/pkg/handler/auth"
	"github.com/erviangelar/go-users-api/pkg/handler/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const defaultUser = `
INSERT INTO public.users
(created_at, username, "name", "role", "password")
VALUES(NOW(), 'admin', 'admin', 'ADMIN', '$2a$14$AFxmldcPxurFLsay/fNtE.NPXUmgregh1VsNUHBmbuLe1m3wby9pO') ON CONFLICT ("username") DO NOTHING;
INSERT INTO public.users
(created_at, username, "name", "role", "password")
VALUES(NOW(), 'user', 'user', 'USER', '$2a$14$pGJf5uGp6F8jTkhYspfoUe4hAGfDgGfbz99KFwJ9Xv8JKtVE1eXpO') ON CONFLICT ("username") DO NOTHING;
`

// @title GO User API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs := config.NewConfigurations()
	// fmt.Println(configs)
	db.Init(configs)
	db, err := db.NewConnection(configs)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.MustExec(defaultUser)

	//Set Api Route Based Group
	r := gin.Default()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("OPTIONS")
	r.Use(cors.New(corsConfig))
	api := r.Group("/api")
	{
		api.GET("/healthcheck", HealthCheck)
		auth.RegisterRoutes(api, db)
		users.RegisterRoutes(api, db)
	}

	url := ginSwagger.URL("http://localhost:3000/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//Start Server
	if err := r.Run(configs.Port); err != nil {
		log.Fatal(err)
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Healthcheck
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/healthcheck [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}
	c.JSON(http.StatusOK, res)
}
