// Copyright Red Hat

package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time	     `json:"updated_at_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`	
}
type AuthRealm struct  {
	BaseModel
	Name		string	`json:"name" gorm:"index:idx_name,unique"`	// Enforce uniqueness on the name within an account
	Account 	string  `json:"account" gorm:"index:idx_name,unique"` // Enforce uniqueness on the name within an account
	CustomResource	datatypes.JSON `json:"custom_resource,omitempty"`
}
