// Copyright Red Hat

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
	"github.com/identitatem/idp-configs-api/pkg/errors"
	"github.com/identitatem/idp-configs-api/pkg/models"
	"github.com/onsi/gomega"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"gorm.io/datatypes"
)

var	(
	authRealm models.AuthRealm
	authRealms []models.AuthRealm
	xrhid identity.XRHID
	badRequest errors.BadRequest
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

	// Add AccountNumber to request context
	ctx := getContextWithAccountNumber("1234567", req)
	GetAuthRealmsForAccount(responseRecorder, req.WithContext(ctx))

	// No records found with Account number 1234567
	json.NewDecoder(responseRecorder.Body).Decode(&authRealms)
	g.Expect(len(authRealms)).To(gomega.Equal(0))

	// Change AccountNumber in req context to match record in DB
	ctx = getContextWithAccountNumber("0000000", req)
	GetAuthRealmsForAccount(responseRecorder, req.WithContext(ctx))
	json.NewDecoder(responseRecorder.Body).Decode(&authRealms)
	g.Expect(len(authRealms)).To(gomega.Equal(1))
	g.Expect(authRealms[0].Name).To(gomega.Equal("TestRecord1"))
}

func TestCreateAuthRealmForAccount (t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	authRealm = models.AuthRealm{
		Account: "1234567",
	}
	authRealmJSON, _ := json.Marshal(&authRealm)

	responseRecorder := httptest.NewRecorder()

	// Expect Bad Request if Name is missing in request body
	req, err := http.NewRequest("POST", "/api/idp-configs-api/v0/auth_realms", bytes.NewBuffer(authRealmJSON))
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	ctx := getContextWithAccountNumber("1234567", req)
	CreateAuthRealmForAccount(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))
	
	// Expect Bad Request if CustomResource is missing in request body
	authRealm.Name = "TestRecord22"
	authRealm.CustomResource = nil
	authRealmJSON, _ = json.Marshal(authRealm)
	responseRecorder = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/api/idp-configs-api/v0/auth_realms", bytes.NewBuffer(authRealmJSON))
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	ctx = getContextWithAccountNumber("1234567", req)
	CreateAuthRealmForAccount(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))	
	json.NewDecoder(responseRecorder.Body).Decode(&badRequest)
	g.Expect(badRequest.Title).To(gomega.Equal("The request body must contain 'Name' and 'CustomResource'"))	

	// Expect 409 if Name is not unique within Account
	authRealm.Name = "TestRecord1" 	// Same as existing record with AccountNumber 0000000
	authRealm.Account = "0000000"
	authRealm.CustomResource = datatypes.JSON([]byte(`{"apiVersion": "identityconfig.identitatem.io/v1alpha1", "kind": "authrealm", "metadata": {"name": "122-auth-realm"}}`))		
	authRealmJSON, _ = json.Marshal(authRealm)
	responseRecorder = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/api/idp-configs-api/v0/auth_realms", bytes.NewBuffer(authRealmJSON))
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	ctx = getContextWithAccountNumber("0000000", req)	
	CreateAuthRealmForAccount(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusConflict))

	// Expect Bad Request if Account number does not match request body
	authRealm.Account = "0000000"	
	authRealmJSON, _ = json.Marshal(authRealm)
	responseRecorder = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/api/idp-configs-api/v0/auth_realms", bytes.NewReader(authRealmJSON))
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	ctx = getContextWithAccountNumber("1234567", req)	// Set different account in request context
	CreateAuthRealmForAccount(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))	

	// Successful response - Name (unique) and CustomResource present in req body, account number taken from request context
	authRealm.Account = ""
	authRealm.Name = "TestRecord3"
	authRealmJSON, _ = json.Marshal(authRealm)
	responseRecorder = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/api/idp-configs-api/v0/auth_realms", bytes.NewReader(authRealmJSON))
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	ctx = getContextWithAccountNumber("1234567", req)	// Set different account in request context
	CreateAuthRealmForAccount(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusOK))
	json.NewDecoder(responseRecorder.Body).Decode(&authRealm)
	g.Expect(authRealm.ID).To(gomega.Equal(uint(0)))
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
	db.DB.Where("Account = ?", "1234567").Unscoped().Delete(&authRealm)
}

func getContextWithAccountNumber(accountNumber string, req *http.Request) context.Context {
	xrhid.Identity.AccountNumber = accountNumber
	ctx := context.WithValue(req.Context(), identity.Key, xrhid)
	
	return ctx
}
