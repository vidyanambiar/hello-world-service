// Copyright Red Hat

package models

import (
	"time"

	"gorm.io/gorm"
)

type AuthRealm struct  {
	gorm.Model
	Id			int		`json:"id"`
	Name		string	`json:"name"`
	Account 	string  `json:"account"` 
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
	CustomResource	map[string]interface{} `json:"customResource"`
}