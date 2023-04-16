package db

import (
	"errors"

	"github.com/stonear/api_gateway/model"
	"gorm.io/gorm"
)

func (db *Db) Migrate() {
	// Migrate the schema
	db.AutoMigrate(&model.Route{})
	if err := db.AutoMigrate(&model.Route{}); err == nil && db.Migrator().HasTable(&model.Route{}) {
		if err := db.First(&model.Route{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Insert seed data
			r := model.Route{
				Name:          "Pokemon",
				Description:   "-",
				Method:        "GET",
				RequestPath:   "/v2/pokemon",
				TargetService: "https://pokeapi.co/api/v2",
				TargetPath:    "pokemon",
				NeedAuth:      false,
			}
			db.Create(&r)
		}
	}
}
