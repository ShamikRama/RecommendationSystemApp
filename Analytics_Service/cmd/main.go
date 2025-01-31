package main

import "Analytics_Service/internal/app"

// @title Product and User Stats API
// @version 1.0
// @description Это API для получения действий пользователя, событий и статистики продуктов.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8085
// @BasePath /
func main() {
	app.RunApp()
}
