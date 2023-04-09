package db

import (
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"wer/cmd/game/helpers"
	"wer/cmd/game/src/structures"
)

type Repository struct {
	db *gorm.DB
}

func Connect() *Repository {

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        helpers.Env("POSTGRESQL_DSN", ""),
	}))
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&structures.Story{}, &structures.StoryLine{}, &structures.StoryLineChoice{}, &structures.Profile{}, &structures.ProfileProgress{})
	if err != nil {
		return nil
	}

	return &Repository{db}
}
