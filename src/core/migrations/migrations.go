package migrations

import (
	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/schema/common"
	"evidence-hub-be/src/core/schema/evident"
	"evidence-hub-be/src/core/schema/evidentPhoto"
	"evidence-hub-be/src/core/schema/leaderboard"
	"log"
)

func RunMigrations() {
	if config.DB == nil {
		log.Fatal("config.DB is nil")
	}
	err := config.DB.AutoMigrate(
		&common.Roles{},
		&common.Users{},
	)
	if err != nil {
		log.Printf("Cannot Migrate Database =========================== %s", err.Error())
	}
}

func RunEvidentMigrations() {
	err := config.DB.AutoMigrate(
		&evident.Evident{},
	)
	if err != nil {
		log.Printf("Cannot Migrate evident Table =========================== %s", err.Error())
	}
}
func RunEvidentPhotoMigrations() {
	err := config.DB.AutoMigrate(
		&evidentPhoto.EvidentPhoto{},
	)
	if err != nil {
		log.Printf("Cannot Migrate evident photo Table =========================== %s", err.Error())
	}
}
func RunLeaderBoardsMigrations() {
	err := config.DB.AutoMigrate(
		&leaderboard.LeaderBoard{},
	)
	if err != nil {
		log.Printf("Cannot Migrate leaderboard Table =========================== %s", err.Error())
	}
}
