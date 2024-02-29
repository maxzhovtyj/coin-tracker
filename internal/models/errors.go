package models

import "errors"

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrWalletAlreadyExists = errors.New("wallet already exists")
	ErrWalletNotFound      = errors.New("wallet not found")
)
