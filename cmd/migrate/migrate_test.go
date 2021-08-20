// Copyright Red Hat

package main

import (
	"testing"

	"github.com/identitatem/idp-configs-api/config"
	"github.com/identitatem/idp-configs-api/pkg/db"
	"github.com/identitatem/idp-configs-api/pkg/models"
	"github.com/onsi/gomega"
	"gorm.io/datatypes"
)

var	authRealm models.AuthRealm

func TestMigrateSchema(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	migrateSchema()

	// Config initialized
	cfg := config.Get()
	g.Expect(cfg.WebPort).To(gomega.Equal(3000))

	// DB initialized
	authRealm = models.AuthRealm{
		Account: "0000000",
		Name:  "TestRecord1",
	}
	authRealm.CustomResource = datatypes.JSON{}
	result := db.DB.Create(&authRealm)
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())

	tearDown()
}
	
func tearDown() {
	db.DB.Where("Account = ?", "0000000").Unscoped().Delete(&authRealm)
}
