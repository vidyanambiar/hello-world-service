// Copyright Red Hat

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	chi "github.com/go-chi/chi/v5"
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

	// TODO: support filtering and searching by name (query param)

	// Respond with auth realms for the account
	json.NewEncoder(w).Encode(&authRealms)	
}

func CreateAuthRealmForAccount(w http.ResponseWriter, r *http.Request) {
    var authRealm models.AuthRealm

	w.Header().Set("Content-Type", "application/json")

	// Get Account from request context
	account, err := common.GetAccount(r)
	if (err != nil) {        
		errors.RespondWithBadRequest(err.Error(), w)
        return		
	}	

    err = json.NewDecoder(r.Body).Decode(&authRealm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// Request body must contain the auth-realm Name and Custom Resource
	if (authRealm.Name == "" || authRealm.CustomResource == nil) {		
		errors.RespondWithBadRequest("The request body must contain 'Name' and 'CustomResource'", w)
		return	
	}

	// TODO: Additional validations on Custom Resource/ evaluating checksum

	// If the request body contains an Account number, it should match the requestor's Account number retrieved from the reqeust context
	if (authRealm.Account != "") {
		err = validateAccount(&authRealm, account)
		if (err != nil) {
			errors.RespondWithBadRequest("Account in the request body does not match account for the authenticated user", w)
			return			
		}		
	}	
	// Set account for the auth realm record
	authRealm.Account = account

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
	authRealmJSON, _ := json.Marshal(authRealm)	
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(authRealmJSON))
}

type key int

const authRealmKey key = 0

// AuthRealmCtx is a handler for Auth Realm requests
func AuthRealmCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var authRealm models.AuthRealm
		// Get account from request header
		account, err := common.GetAccount(r)

		if (err != nil) {        
			errors.RespondWithBadRequest(err.Error(), w)
			return		
		}
		if authRealmID := chi.URLParam(r, "id"); authRealmID != "" {
			// Fetch record based on Auth Realm ID
			result := db.DB.First(&authRealm, authRealmID)

			if (result.Error != nil) {
				errors.RespondWitNotFound(result.Error.Error(), w)
				return
			}
			
			if authRealm.Account != "" {
				// Check that the requestor's account matches the account in the DB
				err = validateAccount(&authRealm, account)
				if (err != nil) {        
					errors.RespondWithForbidden("Requestor's account does not match the Auth Realm account", w)
					return		
				}				
			}
			ctx := context.WithValue(r.Context(), authRealmKey, &authRealm)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})	
}

func GetAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authRealm, ok := ctx.Value(authRealmKey).(*models.AuthRealm)
	if !ok {
		http.Error(w, "must pass id", http.StatusBadRequest)
		return
	}

	// Respond with auth realms for the account
	json.NewEncoder(w).Encode(&authRealm)	
}

func UpdateAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "UpdateAuthRealmByID")
}

func DeleteAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteAuthRealmByID")
}

func validateAccount (authRealm *models.AuthRealm, account string) (error) {	
	if (authRealm.Account != "" && authRealm.Account != account) {
		// Account in the request body must match the account of the authenticated user
		return fmt.Errorf("mismatch in account")
	}

	return nil
}