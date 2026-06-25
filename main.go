package main

import (
	"log"
	"os"
	"sync"

	"evidence-hub-be/server"
	"evidence-hub-be/src/core/commands"
	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/config/envs"
	"evidence-hub-be/src/core/constants"
	"evidence-hub-be/src/core/migrations"

	// "evidence-hub-be/src/core/databases/migrations"
	// "evidence-hub-be/src/core/seeder"

	"github.com/joho/godotenv"
)

// @title API Evident Hub V1
// @version 1.0
// @description EVIDENT HUB API.
// @termsOfService http://swagger.io/terms/
// @contact.name Sofware Engineer
// @contact.email gerry.radityaky@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey token
// @in header
// @name Authorization
// @BasePath /
func main() {

	var once sync.Once
	once.Do(createUploadDir)

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := envs.Get()

	config.ConnectDB(cfg) // WAJIB sebelum migrate

	if cfg.DbConfig.DbMigrate {
		migrate()
	}

	if cfg.DbConfig.DbSeeder {
		go seed(cfg)
	}

	srv := server.New(cfg)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}

func seed(cfg *envs.Config) {
	commands.CreateRoleAdmin()
	commands.CreateUserAdmin(cfg)
}

func migrate() {
	migrations.RunMigrations()
	migrations.RunEvidentMigrations()
	migrations.RunEvidentPhotoMigrations()
	migrations.RunLeaderBoardsMigrations()
}

func createUploadDir() {

	log.Println("creating uploads dir ...")

	if _, err := os.Stat(constants.UploadDir); os.IsNotExist(err) {

		if err := os.Mkdir(constants.UploadDir, os.ModePerm); err != nil {
			panic(err)
		}

		if err := os.Mkdir(
			constants.UploadDir+"/images",
			os.ModePerm,
		); err != nil {
			panic(err)
		}
	}
}
