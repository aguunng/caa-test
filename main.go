package main

import (
	"caa-test/cmd"

	_ "github.com/joho/godotenv/autoload"
)

// @title           Customer Service API
// @version         1.0
// @description     This API provides endpoints for managing customer rooms, handling webhooks, and configuring application settings.
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.url    http://example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI Specification
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cmd.Execute()
}
