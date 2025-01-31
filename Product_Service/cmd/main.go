package main

import (
	"Product_Service/internal/app"

	_ "github.com/lib/pq"
)

// @title Product API
// @version 1.0
// @description Это API для управления продуктами и корзиной.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /
func main() {
	app.RunApp()
}
