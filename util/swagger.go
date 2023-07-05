package util

import "github.com/risk-place-angola/backend-risk-place/api"

func init() {
	env = LoadEnv(".env")
}

// @title Risk Place Angola API
// @description This is a sample server Risk Place Angola server.
// @version 1.0.0
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func Swagger() {
	api.SwaggerInfo.Host = env.REMOTEHOST
	api.SwaggerInfo.BasePath = "/"
	api.SwaggerInfo.Schemes = []string{"http", "https"}
}
