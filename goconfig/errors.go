package goconfig

import "errors"

var (
	// ErrUnmarshalling is the error message for an unmarshalling error.
	ErrUnmarshalling = errors.New("error unmarshalling configuration")
	// ErrVariableNotFound is the error message for a missing environment variable.
	ErrVariableNotFound = errors.New("environment variable not found")
	// ErrUnsupportedExt is the error message for an unsupported extension.
	ErrUnsupportedExt = errors.New("unsupported extension")
	// ErrReadingFile is the error message for a file reading error.
	ErrReadingFile = errors.New("error reading file")
	// ErrOpenDir is the error message for a directory opening error.
	ErrOpenDir = errors.New("error opening directory")
	// ErrOpeningEnvFile is the error message for a configuration reading error.
	ErrOpeningEnvFile = errors.New("error opening .env file")
	// ErrInvalidEnvFormat is the error message for an invalid .env format.
	ErrInvalidEnvFormat = errors.New("invalid .env format")
)
