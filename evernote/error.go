package evernote

import "errors"

var (
	// ErrNotLoggedIn is returned when the user is trying to perform
	// authenticated actions without being authenticated.
	ErrNotLoggedIn = errors.New("your are not logged in")
	// ErrAlreadyLoggedIn is returned if the user is trying to authenticate
	// but is already authenticated.
	ErrAlreadyLoggedIn = errors.New("you are already logged in")
	// ErrTempTokenMismatch is returned if the callback doesn't match the
	// expected token.
	ErrTempTokenMismatch = errors.New("temporary token mismatch")
	// ErrAccessRevoked is returned if the user decline access.
	ErrAccessRevoked = errors.New("access revoked")
	// ErrNoGUIDSet is returned if the note does not have a GUID.
	ErrNoGUIDSet = errors.New("no GUID set.")
	// ErrNoTitleSet is returned if the not does not have a title.
	ErrNoTitleSet = errors.New("no title set")
)
