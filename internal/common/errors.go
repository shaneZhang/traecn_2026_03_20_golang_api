// Copyright 2024 The Refactored Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNotFound            = errors.New("resource not found")
	ErrAlreadyExists       = errors.New("resource already exists")
	ErrInvalidInput        = errors.New("invalid input")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInternalServer      = errors.New("internal server error")
	ErrDatabase            = errors.New("database error")
	ErrValidation          = errors.New("validation error")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTokenExpired        = errors.New("token expired")
	ErrTokenInvalid        = errors.New("token invalid")
	ErrRateLimited         = errors.New("rate limited")
	ErrServiceUnavailable  = errors.New("service unavailable")
	ErrNotImplemented      = errors.New("not implemented")
	ErrBadRequest          = errors.New("bad request")
	ErrConflict            = errors.New("conflict")
)

// AppError represents application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates new application error
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapError wraps an error with context
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// IsNotFound checks if error is not found
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists checks if error is already exists
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsUnauthorized checks if error is unauthorized
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden checks if error is forbidden
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsValidation checks if error is validation error
func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}

// Error codes
const (
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeAlreadyExists      = "ALREADY_EXISTS"
	ErrCodeInvalidInput       = "INVALID_INPUT"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabase           = "DATABASE_ERROR"
	ErrCodeValidation         = "VALIDATION_ERROR"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid       = "TOKEN_INVALID"
	ErrCodeRateLimited        = "RATE_LIMITED"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrCodeNotImplemented     = "NOT_IMPLEMENTED"
	ErrCodeBadRequest         = "BAD_REQUEST"
	ErrCodeConflict           = "CONFLICT"
)

// User errors
var (
	ErrUserNotFound        = NewAppError(ErrCodeNotFound, "user not found", ErrNotFound)
	ErrUserAlreadyExists   = NewAppError(ErrCodeAlreadyExists, "user already exists", ErrAlreadyExists)
	ErrInvalidUserInput    = NewAppError(ErrCodeInvalidInput, "invalid user input", ErrInvalidInput)
	ErrInvalidPassword     = NewAppError(ErrCodeInvalidCredentials, "invalid password", ErrInvalidCredentials)
	ErrUserDisabled        = NewAppError(ErrCodeForbidden, "user is disabled", ErrForbidden)
	ErrUserNotVerified     = NewAppError(ErrCodeForbidden, "user is not verified", ErrForbidden)
)

// Organization errors
var (
	ErrOrganizationNotFound      = NewAppError(ErrCodeNotFound, "organization not found", ErrNotFound)
	ErrOrganizationAlreadyExists = NewAppError(ErrCodeAlreadyExists, "organization already exists", ErrAlreadyExists)
	ErrOrganizationHasUsers      = NewAppError(ErrCodeConflict, "organization has users", ErrConflict)
)

// Application errors
var (
	ErrApplicationNotFound      = NewAppError(ErrCodeNotFound, "application not found", ErrNotFound)
	ErrApplicationAlreadyExists = NewAppError(ErrCodeAlreadyExists, "application already exists", ErrAlreadyExists)
	ErrInvalidClientCredentials = NewAppError(ErrCodeInvalidCredentials, "invalid client credentials", ErrInvalidCredentials)
	ErrInvalidRedirectURI       = NewAppError(ErrCodeInvalidInput, "invalid redirect URI", ErrInvalidInput)
	ErrInvalidGrantType         = NewAppError(ErrCodeInvalidInput, "invalid grant type", ErrInvalidInput)
	ErrInvalidScope             = NewAppError(ErrCodeInvalidInput, "invalid scope", ErrInvalidInput)
)

// MFA errors
var (
	ErrMFANotEnabled        = NewAppError(ErrCodeForbidden, "MFA not enabled", ErrForbidden)
	ErrInvalidMFACode       = NewAppError(ErrCodeInvalidCredentials, "invalid MFA code", ErrInvalidCredentials)
	ErrMFAAlreadyEnabled    = NewAppError(ErrCodeConflict, "MFA already enabled", ErrConflict)
	ErrMFASetupRequired     = NewAppError(ErrCodeForbidden, "MFA setup required", ErrForbidden)
	ErrRecoveryCodeInvalid  = NewAppError(ErrCodeInvalidCredentials, "invalid recovery code", ErrInvalidCredentials)
	ErrRecoveryCodeUsed     = NewAppError(ErrCodeInvalidCredentials, "recovery code already used", ErrInvalidCredentials)
)
