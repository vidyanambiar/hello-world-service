package services

import (
	"fmt"
	"net/http"
)

func GetAuthRealmsForAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GetAuthRealmsForAccount")
}

func CreateAuthRealmForAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "CreateAuthRealmForAccount")
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