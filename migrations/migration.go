package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/ponyjackal/eventual-consistency-blog/infra/database"
	"github.com/ponyjackal/eventual-consistency-blog/models"
)

func Migrate() {
	m := gormigrate.New(database.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// Define your migration steps here
		{
			ID: "20230705120000",
			Migrate: func(tx *gorm.DB) error {
				// Perform migration operations to create the tables and modify the schema
				return tx.AutoMigrate(
					&models.Post{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				// Define rollback operations if necessary
				return tx.Migrator().DropTable(
					&models.Post{},
				)
			},
		},
	})
	
	// Run the migrations
	if err := m.Migrate(); err != nil {
		panic(err)
	}
}
