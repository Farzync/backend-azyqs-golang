package errors

import "errors"

var (
	ErrUsernameTaken    = errors.New("username_already_taken")
	ErrEmailTaken       = errors.New("email_already_taken")
	ErrUserNotFound     = errors.New("user_not_found")
	ErrInvalidPassword  = errors.New("invalid_password")
	ErrPasswordHash     = errors.New("password_hash_error")
	ErrDuplicateRecord  = errors.New("duplicate_record")
	ErrUserDeleteFailed = errors.New("user_delete_failed")
	ErrUserUpdateFailed = errors.New("user_update_failed")
	ErrPasswordMismatch = errors.New("password_mismatch")
	ErrInvalidInput     = errors.New("invalid_input")
	ErrUserIDNotFound   = errors.New("user_id_not_found")
	ErrInvalidUserID   = errors.New("user_id_is_invalid")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInternalServer   = errors.New("internal_server_error")
)
