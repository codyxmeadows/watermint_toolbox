package dbx_error

import (
	"strings"
)

type DropboxError struct {
	ErrorTag          string `path:"error.\\.tag" json:"error,omitempty"`
	ErrorSummary      string `path:"error_summary" json:"error_summary,omitempty"`
	UserMessageLocale string `json:"user_message_lang,omitempty"`
	UserMessage       string `json:"user_message,omitempty"`
}

func (z DropboxError) Error() string {
	return z.ErrorSummary
}

func (z DropboxError) HasPrefix(prefix string) bool {
	return strings.HasPrefix(z.ErrorSummary, prefix)
}

func NewErrors(err error) Errors {
	if err == nil {
		return nil
	}
	if de, ok := err.(*DropboxError); ok {
		return &errorsImpl{de: *de}
	} else {
		return &errorsImpl{de: DropboxError{}}
	}
}

type Errors interface {
	Auth() ErrorAuth
	Access() ErrorAccess
	Path() ErrorEndpointPath
	Endpoint() ErrorEndpoint
	To() ErrorWrite
	BadPath() ErrorSharePath

	// too_many_write_operations
	IsTooManyWriteOperations() bool
	// too_many_files
	IsTooManyFiles() bool
	// folder_name_already_used
	IsFolderNameAlreadyUsed() bool
	// group_name_already_used
	IsGroupNameAlreadyUsed() bool

	Summary() string
}

// 401: Bad or expired token. This can happen if the access token is expired or if the access token has been revoked by Dropbox or the user.
type ErrorAuth interface {
	// The access token is invalid.
	IsInvalidAccessToken() bool
	// The user specified in 'Dropbox-API-Select-User' is no longer on the team.
	IsInvalidSelectUser() bool
	// The user specified in 'Dropbox-API-Select-Admin' is not a Dropbox Business team admin.
	IsInvalidSelectAdmin() bool
	// The user has been suspended.
	IsUserSuspended() bool
	// The access token has expired.
	IsExpiredAccessToken() bool
	// The access token does not have the required scope to access the route.
	IsMissingScope() bool
	// The route is not available to public.
	IsRouteAccessDenied() bool
}

// 403: The user or team account doesn't have access to the endpoint or feature.
type ErrorAccess interface {
	// Current account type cannot access the resource.
	IsInvalidAccountType() bool
	// Current account cannot access Paper.
	IsPaperAccessDenied() bool
}

// 409: Endpoint specific error, Path
type ErrorEndpointPath interface {
	// There is nothing at the given path.
	IsNotFound() bool
	// The given path does not satisfy the required path format
	IsMalformedPath() bool
	// Couldn't write to the target path because there was something in the way.
	IsConflict() bool
	// There are too many write operations in user's Dropbox. Please retry this request.
	IsTooManyWriteOperations() bool
}

// 409: Other endpoint specific error
type ErrorEndpoint interface {
	IsRateLimit() bool
}

// 409: WriteError
type ErrorWrite interface {
	IsConflict() bool
}

// 409: SharePathError
type ErrorSharePath interface {
	IsAlreadyShared() bool
}

type errorsImpl struct {
	de DropboxError
}

func (z errorsImpl) IsGroupNameAlreadyUsed() bool {
	return z.de.HasPrefix("group_name_already_used")
}

func (z errorsImpl) BadPath() ErrorSharePath {
	return NewSharePath("bad_path", z.de)
}

func (z errorsImpl) IsFolderNameAlreadyUsed() bool {
	return z.de.HasPrefix("folder_name_already_used")
}

func (z errorsImpl) IsTooManyFiles() bool {
	return z.de.HasPrefix("too_many_files")
}

func (z errorsImpl) IsTooManyWriteOperations() bool {
	return z.de.HasPrefix("too_many_write_operations") || z.Path().IsTooManyWriteOperations()
}

func (z errorsImpl) Summary() string {
	return z.de.ErrorSummary
}

func (z errorsImpl) Auth() ErrorAuth {
	return NewErrorAuth(z.de)
}

func (z errorsImpl) Access() ErrorAccess {
	return NewErrorAccess(z.de)
}

func (z errorsImpl) Path() ErrorEndpointPath {
	return NewErrorPath(z.de)
}

func (z errorsImpl) Endpoint() ErrorEndpoint {
	return NewErrorEndpoint(z.de)
}

func (z errorsImpl) To() ErrorWrite {
	return NewErrorWrite("to", z.de)
}
