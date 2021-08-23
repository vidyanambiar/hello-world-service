// Copyright Red Hat

package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)
type AuthRealm struct  {
	gorm.Model
	Name		string	`json:"Name" gorm:"index:idx_name,unique"`	// Enforce uniqueness on the name within an account
	Account 	string  `json:"Account" gorm:"index:idx_name,unique"` // Enforce uniqueness on the name within an account
	CustomResource	datatypes.JSON `json:"CustomResource,omitempty"`
}