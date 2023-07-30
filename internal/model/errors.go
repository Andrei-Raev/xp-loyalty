package model

import "errors"

var (
	ErrUserExists              = errors.New("user exists")
	ErrUserNotFound            = errors.New("My user not found")
	ErrInvalidAccessToken      = errors.New("access token in invalid")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrWrongPassword           = errors.New("worng password")
	ErrWrongRole               = errors.New("this role has no access")
	ErrNoSuchCard              = errors.New("no such card")
	ErrCardAlreadyInPool       = errors.New("card already exist in pool")
	ErrInterfaceCast           = errors.New("couldn't cast interface")
	ErrNoSuchPool              = errors.New("no such pool")
	ErrSmthWrong               = errors.New("something went wrong")
	ErrNoRandomCards           = errors.New("random cards slice in empty")
)
