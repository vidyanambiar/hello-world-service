// Copyright Red Hat

package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/identitatem/idp-configs-api/pkg/common"
	"github.com/identitatem/idp-configs-api/pkg/db"
	"github.com/identitatem/idp-configs-api/pkg/errors"
	"github.com/identitatem/idp-configs-api/pkg/models"
)

func GetAuthRealmsForAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authRealms []models.AuthRealm

	// Get account from request header
	account, err := common.GetAccount(r)
	fmt.Println("&&& account", account)

	if (err != nil) {        
		errors.RespondWithBadRequest(err.Error(), w)
        return		
	}
	
	// Fetch Auth Realms for specific account from the DB
	result := db.DB.Where("Account = ?", account).Find(&authRealms)

	if result.Error != nil {		
		errors.RespondWithBadRequest(result.Error.Error(), w)
		return
	}

	fmt.Println("&&& authRealms", authRealms)

	// TODO: support filtering and searching by name (query param)

	// Respond with auth realms for the account
	json.NewEncoder(w).Encode(&authRealms)	
}

func CreateAuthRealmForAccount(w http.ResponseWriter, r *http.Request) {
    var authRealm models.AuthRealm

	w.Header().Set("Content-Type", "application/json")

    err := json.NewDecoder(r.Body).Decode(&authRealm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// Request body must contain the auth-realm Name and Custom Resource
	if (authRealm.Name == "" || authRealm.CustomResource == nil) {
		// Bad request
		errors.RespondWithBadRequest("The request body must contain 'Name' and 'CustomResource'", w)
		return	
	}

	// TODO: Additional validations on Custom Resource/ evaluating checksum

	// Validate account from request
	err = checkAccount(&authRealm, r)

	if (err != nil) {
		// Bad request
		errors.RespondWithBadRequest(err.Error(), w)
		return			
	}

	// Temporarily responding with the auth realm object that will be submitted to the DB
	authRealmJSON, _ := json.Marshal(authRealm)

	// Create record for auth realm in the DB		
	tx := db.DB.Create(&authRealm)
	if tx.Error != nil {
		errorMessage := tx.Error.Error()
		if (strings.Contains(strings.ToLower(errorMessage), "unique constraint")) {	// The error message looks a little different between sqlite and postgres 
			// Unique constraint violated (return 409)
			errors.RespondWithConflict("Error creating record in the DB: " + tx.Error.Error(), w)			
		} else {
			// Error updating the DB		
			errors.RespondWithInternalServerError("Error creating record in the DB: " + tx.Error.Error(), w)
		}
		return			
	}

	// Return ID for created record (Temporily responding with the complete record)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(authRealmJSON))
}

func GetAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GetAuthRealmByID")
}

func UpdateAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "UpdateAuthRealmByID")
}

func DeleteAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteAuthRealmByID")
}

func checkAccount (authRealm *models.AuthRealm, r *http.Request) (error) {	
	// Get account from request header
	account, err := common.GetAccount(r)

	if err != nil {
		return err
	}

	// Check if request body contains an Account (optional field)
	if (authRealm.Account != "" && authRealm.Account != account) {
		// Account in the request body must match the account of the authenticated user
		return fmt.Errorf("account in the request body does not match account for the authenticated user")
	}

	// Set the account from the request header as the account for the auth realm
	authRealm.Account = account

	return nil
}