// Copyright Red Hat

package common

import (
	"fmt"
	"net/http"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// Get account from http request header
func GetAccount(r *http.Request) (string, error) {
	if r.Context().Value(identity.Key) != nil {
		ident := identity.Get(r.Context())
		fmt.Println("ident: ", ident)
		fmt.Println("ident.Identity: ", ident.Identity)
		fmt.Println("ident.Identity.AccountNumber: ", ident.Identity.AccountNumber)
		if ident.Identity.AccountNumber != "" {
			return ident.Identity.AccountNumber, nil
		}
	}
	return "", fmt.Errorf("cannot find account number")

}