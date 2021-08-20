// // Copyright Red Hat

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/identitatem/idp-configs-api/config"
	"github.com/identitatem/idp-configs-api/pkg/db"
	"github.com/identitatem/idp-configs-api/pkg/models"
	"github.com/onsi/gomega"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"gorm.io/datatypes"
)

var	(
	authRealm models.AuthRealm
	authRealms []models.AuthRealm
	xrhid identity.XRHID
)

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func TestGetAuthRealmsForAccount (t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	
	responseRecorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/idp-configs-api/v0/auth_realms", nil)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	// Bad request if there is no AccountNumber present in the request context
	GetAuthRealmsForAccount(responseRecorder, req)
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))

	// Mock AccountNumber for request context
	xrhid = identity.XRHID{
		Identity: identity.Identity{
			AccountNumber: "1234567",
		},
	}
	// Add AccountNumber to request context
	ctx := context.WithValue(req.Context(), identity.Key, xrhid)

	GetAuthRealmsForAccount(responseRecorder, req.WithContext(ctx))

	// No records found with Account number 1234567
	json.NewDecoder(responseRecorder.Body).Decode(&authRealms)
	g.Expect(len(authRealms)).To(gomega.Equal(0))

	// Change AccountNumber in req context to match record in DB
	xrhid.Identity.AccountNumber = "0000000"
	ctx = context.WithValue(req.Context(), identity.Key, xrhid)
	GetAuthRealmsForAccount(responseRecorder, req.WithContext(ctx))
	json.NewDecoder(responseRecorder.Body).Decode(&authRealms)
	g.Expect(len(authRealms)).To(gomega.Equal(1))
	g.Expect(authRealms[0].Name).To(gomega.Equal("TestRecord1"))
}

func TestCreateAuthRealmForAccount (t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	authRealm = models.AuthRealm{
		Account: "0000000",
		Name:  "TestRecord2",
	}
	authRealm.CustomResource = datatypes.JSON([]byte(`{
        "apiVersion": "identityconfig.identitatem.io/v1alpha1",
        "kind": "authrealm",
        "metadata": {
            "name": "122-auth-realm"
		}
	}`))	
	authRealmJSON, _ := json.Marshal(authRealm)

	responseRecorder := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/idp-configs-api/v0/auth_realms", bytes.NewReader(authRealmJSON))
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	// Bad request if Name
	GetAuthRealmsForAccount(responseRecorder, req)
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))

	// Mock AccountNumber for request context
	xrhid = identity.XRHID{
		Identity: identity.Identity{
			AccountNumber: "1234567",
		},
	}
	// Add AccountNumber to request context
	ctx := context.WithValue(req.Context(), identity.Key, xrhid)

	GetAuthRealmsForAccount(responseRecorder, req.WithContext(ctx))

	// No records found with Account number 1234567
	json.NewDecoder(responseRecorder.Body).Decode(&authRealms)
	g.Expect(len(authRealms)).To(gomega.Equal(0))

	// Change AccountNumber in req context to match record in DB
	xrhid.Identity.AccountNumber = "0000000"
	ctx = context.WithValue(req.Context(), identity.Key, xrhid)
	GetAuthRealmsForAccount(responseRecorder, req.WithContext(ctx))
	json.NewDecoder(responseRecorder.Body).Decode(&authRealms)
	g.Expect(len(authRealms)).To(gomega.Equal(1))
	g.Expect(authRealms[0].Name).To(gomega.Equal("TestRecord1"))
}

func setUp() {
	// Initialize config for test
	config.Init()
	db.InitDB()

	// Add a record to the DB
	authRealm = models.AuthRealm{
		Account: "0000000",
		Name:  "TestRecord1",
	}
	authRealm.CustomResource = datatypes.JSON{}
	db.DB.Create(&authRealm)	
}

func tearDown() {
	db.DB.Where("Account = ?", "0000000").Unscoped().Delete(&authRealm)
}

func setAccountNumberInRequestContext (accountNumber string) {
	// Mock AccountNumber for request context
	xrhid = identity.XRHID{
		Identity: identity.Identity{
			AccountNumber: accountNumber,
		},
	}	
}