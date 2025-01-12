package main

import (
	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/auth/internal/db"
	"github.com/juxue97/common"
)

//	@title			Auth API
//	@description	This is a authentication backend server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	// Consume .env here
	cfg := config.Configs
	defer db.PgDB.Close()
	defer common.Logger.Sync()

	app := &application{
		config: cfg,
	}
	mux := app.mount()

	common.Logger.Fatal(app.run(mux))
}
