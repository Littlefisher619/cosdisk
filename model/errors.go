package model

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	ErrUnimplemented       = errors.New("unimplemented")
	// FileNotFound will throw if the requested file is not exists
	ErrFileNotFound = errors.New("Your requested file is not found")
	ErrFileExists   = errors.New("Your requested file is already exists")
	// ErrInvalidFileKey will throw if Invalid userfile key
	ErrInvalidFileKey = errors.New("Invalid userfile key")
	// ErrIsDir will throw if This file is a dir
	ErrIsDir = errors.New("This file is a dir")
	// ErrIsFile will throw if This dir is a file
	ErrIsFile          = errors.New("This dir is a file")
	ErrUserNotFound    = errors.New("User not found")
	ErrShareIdNotFound = errors.New("Share id not found")
	ErrShareExpired    = errors.New("The share id is expired")
)
