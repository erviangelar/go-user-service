package config

import (
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

// Configurations wraps all the config variables required by the auth service
type Configurations struct {
	Port                       string
	ServerAddress              string
	DBConn                     string
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	JwtExpiration              int
	TokenHash                  string
}

// NewConfigurations returns a new Configuration object
func NewConfigurations() *Configurations {

	viper.SetConfigFile("./envs/.env")
	viper.ReadInConfig()

	dbURL := viper.GetString("DATABASE_URL")
	conn, _ := pq.ParseURL(dbURL)
	// fmt.Println(dbURL)
	// fmt.Println(conn)

	configs := &Configurations{
		Port:                       viper.GetString("PORT"),
		ServerAddress:              viper.GetString("SERVER_ADDRESS"),
		DBConn:                     conn,
		JwtExpiration:              viper.GetInt("JWT_EXPIRATION"),
		AccessTokenPrivateKeyPath:  viper.GetString("ACCESS_TOKEN_PRIVATE_KEY_PATH"),
		AccessTokenPublicKeyPath:   viper.GetString("ACCESS_TOKEN_PUBLIC_KEY_PATH"),
		RefreshTokenPrivateKeyPath: viper.GetString("REFRESH_TOKEN_PRIVATE_KEY_PATH"),
		RefreshTokenPublicKeyPath:  viper.GetString("REFRESH_TOKEN_PUBLIC_KEY_PATH"),
		TokenHash:                  viper.GetString("TokenHash"),
	}

	// reading heroku provided port to handle deployment with heroku
	port := viper.GetString("PORT")
	if port != "" {
		configs.ServerAddress = "0.0.0.0:" + port
	}
	return configs
}
