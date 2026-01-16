// Package errs defines application-specific error values
// for validation, cache, and internal server errors.
package errs

import "errors"

var (
	ErrInvalidJSON    = errors.New("invalid JSON format")                         // invalid JSON format
	ErrInternal       = errors.New("internal server error")                       // internal server error
	ErrTitleEmpty     = errors.New("chat title cannot be empty")                  // chat title cannot be empty
	ErrTitleTooLong   = errors.New("chat title exceeds maximum length")           // chat title exceeds maximum length
	ErrMessageEmpty   = errors.New("message text cannot be empty")                // message text cannot be empty
	ErrMessageTooLong = errors.New("message text exceeds maximum length")         // message text exceeds maximum length
	ErrInvalidChatID  = errors.New("invalid chat ID; must be a positive integer") // invalid chat ID; must be a positive integer
	ErrChatNotFound   = errors.New("chat not found")                              // chat not found
	ErrLimitTooSmall  = errors.New("limit cannot be negative")                    // limit cannot be negative
	ErrLimitTooLarge  = errors.New("number of messages exceeds service limit")    // number of messages exceeds service limit
	ErrInvalidLimit   = errors.New("invalid limit; must be an integer")           // invalid limit; must be a positive integer
	ErrCacheMiss      = errors.New("cache miss")                                  // cache miss
)
