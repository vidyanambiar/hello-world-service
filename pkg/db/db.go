// Copyright Red Hat

package db

import (
	"fmt"

	"github.com/identitatem/idp-configs-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB ORM variable
var DB *gorm.DB

// InitDB to configure database connectivity
func InitDB() {
	var err error
	var dia gorm.Dialector
	cfg := config.Get()

	if cfg.Database != nil {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			cfg.Database.Hostname,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.Port,
		)
		dia = postgres.Open(dsn)
		fmt.Println("Opening postgres DB...")
	} else {
		dia = sqlite.Open("test.db")
		fmt.Println("Opening sqlite DB...")
	}

	DB, err = gorm.Open(dia, &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err.Error()))
	}
}
