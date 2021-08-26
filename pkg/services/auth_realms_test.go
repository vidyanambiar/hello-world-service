// Copyright Red Hat

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	chi "github.com/go-chi/chi/v5"
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
	g.Expect(len(authRealms) > 0).To(gomega.Equal(true))
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
	g.Expect(badRequest.Title).To(gomega.Equal("The request body must contain 'name' and 'custom_resource'"))	

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

	// Successful response - name (unique) and custom_resource present in req body, account number taken from request context
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
	g.Expect(authRealm.ID).ToNot(gomega.BeNil())
}

func TestGetAuthRealmByID (t *testing.T) {
	g := gomega.NewGomegaWithT(t)
		
	// Get a test auth realm from DB
	result := db.DB.First(&authRealm)
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())
	resultName := authRealm.Name

	// Bad Request if auth realm is not found - no auth realm set in request context
	responseRecorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/idp-configs-api/v0/auth_realms/1000", nil)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	GetAuthRealmByID(responseRecorder, req)
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))

	// successful response
	responseRecorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/idp-configs-api/v0/auth_realms/", nil)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	// Add auth realm to the context
	ctx := context.WithValue(req.Context(), AuthRealmKey, &authRealm)
	GetAuthRealmByID(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusOK))
	json.NewDecoder(responseRecorder.Body).Decode(&authRealm)
	g.Expect(authRealm.Name).To(gomega.Equal(resultName))
}

func TestUpdateAuthRealmByID (t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	
	var authRealmFromDB models.AuthRealm

	var authRealmForRequestBody = models.AuthRealm{
		Name: "NewName",
		Account: "1000000",
	}
	authRealmForRequestBody.CustomResource = datatypes.JSON([]byte(`{"apiVersion": "identityconfig.identitatem.io/v1alpha1", "kind": "authrealm", "metadata": {"name": "122-auth-realm"}}`))
	
	// Get a test auth realm from DB
	result := db.DB.First(&authRealmFromDB)
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())
	authRealmTestID := authRealmFromDB.ID	
	originalName := authRealmFromDB.Name
	fmt.Println("originalName:", originalName)

	// Bad Request if auth realm is not found - no auth realm set in request context
	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/idp-configs-api/v0/auth_realms/1", new(bytes.Buffer))
	UpdateAuthRealmByID(responseRecorder, req)
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))

	// Bad Request - PUT request body is empty
	responseRecorder = httptest.NewRecorder()
	// Add auth realm to the context (auth realm ID was provided as a URL param)
	ctx := context.WithValue(req.Context(), AuthRealmKey, &authRealmFromDB)
	UpdateAuthRealmByID(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))

	// Bad Request - Account number in request body does not match auth realm account in DB
	responseRecorder = httptest.NewRecorder()
	authRealmJSON, _ := json.Marshal(&authRealmForRequestBody)
	req, _ = http.NewRequest("PUT", "/api/idp-configs-api/v0/auth_realms/1", bytes.NewBuffer(authRealmJSON))
	UpdateAuthRealmByID(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))	

	// Successful response - change auth realm name
	responseRecorder = httptest.NewRecorder()
	authRealmForRequestBody.Account = authRealmFromDB.Account
	authRealmJSON, _ = json.Marshal(authRealmForRequestBody)
	req, _ = http.NewRequest("PUT", "/api/idp-configs-api/v0/auth_realms/1", bytes.NewBuffer(authRealmJSON))
	UpdateAuthRealmByID(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusOK))

	// Read auth realm from DB to verify name change
	result = db.DB.First(&authRealmFromDB, authRealmTestID)
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())
	g.Expect(authRealmFromDB.Name).To(gomega.Equal("NewName"))

	// Reset name
	authRealmFromDB.Name = originalName
	db.DB.Save(&authRealmFromDB)
}

func TestDeleteAuthRealmByID (t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// Add a record to the DB
	var authRealmTest = models.AuthRealm{
		Account: "0000000",
		Name:  "TestRecordForDeletion",
	}
	authRealmTest.CustomResource = datatypes.JSON{}
	db.DB.Create(&authRealmTest)

	// Bad Request if auth realm is not found - no auth realm set in request context
	responseRecorder := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/api/idp-configs-api/v0/auth_realms/1000", nil)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	DeleteAuthRealmByID(responseRecorder, req)
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))

	// Successful response
	responseRecorder = httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), AuthRealmKey, &authRealmTest)
	DeleteAuthRealmByID(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusOK))

	// Verify that this record is not found in the DB
	result := db.DB.First(&authRealm, "Name = ?", "TestRecordForDeletion")
	g.Expect(result.Error).Should(gomega.HaveOccurred())

	// Internal Server Error if there is a problem with deleting the record
	responseRecorder = httptest.NewRecorder()
	var nonExistentRecord = models.AuthRealm{
		Account: "0000000",
		Name:  "NonExistent",
	}
	nonExistentRecord.CustomResource = datatypes.JSON{}	
	ctx = context.WithValue(req.Context(), AuthRealmKey, &nonExistentRecord)
	DeleteAuthRealmByID(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusInternalServerError))
}

func TestAuthRealmCtx (t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	next := http.HandlerFunc(testHandler)
	authRealmHandler := AuthRealmCtx(next)
	g.Expect(authRealmHandler).NotTo(gomega.BeNil())

	// Get ID for an existing auth realm
	result := db.DB.First(&authRealm)
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())
	authRealmTestID := authRealm.ID	

	// Bad request - Account number missing from request context
	responseRecorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/idp-configs-api/v0/auth_realms/1000", nil)	
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	authRealmHandler.ServeHTTP(responseRecorder, req)
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))

	// Not Found - Auth Realm ID not found in the DB
	responseRecorder = httptest.NewRecorder()	
	ctx := getContextWithAccountNumber("0000000", req) // Add account number to context	
	ctx = context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{	// Add auth realm id to context
		URLParams: chi.RouteParams{
			Keys:   []string{"id"},
			Values: []string{"1000"},	// This id does not exist in the DB
		},
	})
	authRealmHandler.ServeHTTP(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusNotFound))

	// 403 Forbidden - Account number from request does not match auth realm account
	responseRecorder = httptest.NewRecorder()
	ctx = getContextWithAccountNumber("AccountThatDoesNotMatch", req)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{	// Add auth realm id to context
		URLParams: chi.RouteParams{
			Keys:   []string{"id"},
			Values: []string{fmt.Sprint(authRealmTestID)},
		},
	})
	authRealmHandler.ServeHTTP(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusForbidden))	

	// Successful response
	responseRecorder = httptest.NewRecorder()
	result = db.DB.First(&authRealm)
	g.Expect(result.Error).ShouldNot(gomega.HaveOccurred())
	authRealmTestID = authRealm.ID	
	authRealmTestAccount := authRealm.Account
	ctx = getContextWithAccountNumber(authRealmTestAccount, req)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{	// Add auth realm id to context
		URLParams: chi.RouteParams{
			Keys:   []string{"id"},
			Values: []string{fmt.Sprint(authRealmTestID)},
		},
	})
	authRealmHandler.ServeHTTP(responseRecorder, req.WithContext(ctx))
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusOK))	
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func setUp() {
	// Initialize config for test
	config.Init()
	db.InitDB()

	// Auto-migrate models
	db.DB.AutoMigrate(&models.AuthRealm{})	

	// Add a record to the DB
	authRealm = models.AuthRealm{
		Account: "0000000",
		Name:  "TestRecord1",
	}
	authRealm.CustomResource = datatypes.JSON{}
	db.DB.Create(&authRealm)
}

func tearDown() {
	result := db.DB.Exec("DELETE FROM auth_realms")
	fmt.Println("result.RowsAffected", result.RowsAffected)
}

func getContextWithAccountNumber(accountNumber string, req *http.Request) context.Context {
	xrhid.Identity.AccountNumber = accountNumber
	ctx := context.WithValue(req.Context(), identity.Key, xrhid)
	
	return ctx
}
