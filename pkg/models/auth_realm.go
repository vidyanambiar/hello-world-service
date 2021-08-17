// Copyright Red Hat

package models

import "time"

type AuthRealm struct  {
	Id			int		`json:"id"`
	Name		string	`json:"name"`
	Account 	string  `json:"account"` 
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
	CustomResource	map[string]interface{} `json:"customResource"`
}