// Copyright Red Hat

package db

import (
	"testing"

	"github.com/identitatem/idp-configs-api/config"
	"github.com/identitatem/idp-configs-api/pkg/models"
	"github.com/onsi/gomega"
	"gorm.io/datatypes"
)

func TestInitDB(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	// Initialize config for test
	config.Init()
	
	InitDB()
	
	authRealm := models.AuthRealm{
		Account: "0000000",
		Name:  "TestRecord",
	}
	authRealm.CustomResource = datatypes.JSON{}
	DB.Create(&authRealm)

	result := DB.First(&authRealm)

	// DB should initialize
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())

	DB.Delete(&authRealm)
}