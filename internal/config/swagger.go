package config

import "github.com/risk-place-angola/backend-risk-place/api"

// @title Risk Place Angola API
// @version 1.0.0
// @description This is the API documentation for the Risk Place Angola application.
// @description
// @description ## Environments
// @description - **Development**: https://risk-place-angola-904a.onrender.com
// @description - **Local**: http://localhost:8000

// @contact.name API Support
// @contact.url http://www.riskplace.ao
// @contact.email support@riskplace.ao

// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func Swagger() {
	api.SwaggerInfo.Schemes = []string{"http", "https"}
}
