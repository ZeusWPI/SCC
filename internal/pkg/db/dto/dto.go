// Package dto provides the data transfer objects for the database
package dto

import "github.com/go-playground/validator/v10"

// Validate is a validator instance for JSON transferable objects
var Validate = validator.New(validator.WithRequiredStructEnabled())
