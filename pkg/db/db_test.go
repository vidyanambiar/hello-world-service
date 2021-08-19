// Copyright Red Hat

package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/identitatem/idp-configs-api/config"
	"github.com/identitatem/idp-configs-api/pkg/models"
	"github.com/onsi/gomega"
	"gorm.io/datatypes"
)

var	authRealm models.AuthRealm
var	authRealms []models.AuthRealm

func TestMain(m *testing.M) {
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func TestInitDB(t *testing.T) {
	fmt.Println("TestInitDB")
	g := gomega.NewGomegaWithT(t)
	// Initialize config for test
	config.Init()
	
	InitDB()

	// DB should initialize correctly (verify by adding and querying a record)
	authRealm = models.AuthRealm{
		Account: "0000000",
		Name:  "TestRecord1",
	}
	authRealm.CustomResource = datatypes.JSON{}
	result := DB.Create(&authRealm)
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())

	result = DB.Find(&authRealms)
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())
	g.Expect(authRealms[0].Account).To(gomega.Equal("0000000"))	
}

func tearDown() {
	DB.Where("Account = ?", "0000000").Unscoped().Delete(&authRealm)
}