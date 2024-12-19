package main

import "onlineShop/internal/app"
import _ "onlineShop/docs"

// @title Online Shop API
// @version 1.0
// @description This is a sample server for an online shop.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	app.Run()
}
