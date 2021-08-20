// Copyright Red Hat

package main

import (
	"github.com/identitatem/idp-configs-api/config"
	l "github.com/identitatem/idp-configs-api/logger"
	"github.com/identitatem/idp-configs-api/pkg/db"
	"github.com/identitatem/idp-configs-api/pkg/models"
	log "github.com/sirupsen/logrus"
)

func main() {
	migrateSchema()
}

func migrateSchema () {
	config.Init()
	l.InitLogger()
	cfg := config.Get()
	log.WithFields(log.Fields{
		"Hostname":           cfg.Hostname,
		"Auth":               cfg.Auth,
		"WebPort":            cfg.WebPort,
		"MetricsPort":        cfg.MetricsPort,
		"LogLevel":           cfg.LogLevel,
		"Debug":              cfg.Debug,
	}).Info("Configuration Values:")	
	db.InitDB()
	err := db.DB.AutoMigrate(&models.AuthRealm{})
	if err != nil {
		panic(err)
	}
	log.Info("Migration Completed")	
}