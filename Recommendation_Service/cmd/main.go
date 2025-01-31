package main

import "Recommendation_Service/internal/app"

// @title Recommendation API
// @version 1.0
// @description Это API для получения рекомендаций для пользователей.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8084
// @BasePath /
func main() {
	app.RunApp()
}
