package logs

const (
	FailedToDecodeRequest = "failed to decode request: %v"
	InternalDBError       = "internal database error: %v"
)

// successful users' logs
const (
	UserRegSuccessfully       = "user was registered successfully with id = %d"
	UserLoginedSuccessfully   = "user was logined successfully (id = %d)"
	UserUpdatedSuccessfully   = "user was updated successfully (id = %d)"
	UserDeletedSuccessfully   = "user was deleted successfully (id = %d)"
	UserGivenByIDSuccessfully = "successfull getting user data by id"
)

// failed user's logs
const (
	FailedToGetUserIDFromCtx = "failed to get userID from context"
	FailedToValidateUser     = "failed to validate user: %v"
	FailedToComparePasswords = "failed to compare passwords: %v"
	UnauthenticatedUser      = "user is unauthenticated"
	FailedToHashPassword     = "failed to hash passwords: %v"
)

// successful password's logs
const (
	PasswordSavedSuccessfully     = "password was saved successfully"
	PasswordUpdatedSuccessfully   = "password was updated successfully"
	PasswordsGivenSuccessfully    = "passwords was given successfully"
	PasswordDeletedSuccessfully   = "passwords was deleted successfully"
	PasswordGeneratedSuccessfully = "password was generated successfully"
)

// failed password's logs
const (
	FailedToGetPasswordIDFromURL     = "failed to get passwordID from url"
	FailedToGetPasswordSettings      = "failed to get password setting for geneating from URL: %v"
	FailedToValidatePassword         = "failed to validate password: %v"
	FailedToValidatePasswordSettings = "failed to validate password settings: %v"
)

// session failed logs
const (
	FailedToGetSession           = "failed to get session"
	FailedToSaveSession          = "failed to save session"
	FailedToGetSessionKey        = "failed to get session key: %v"
	FailedToGetUserIDFromSession = "failed to get user id from session"
)

// main.go logs
const (
	FailedToInitConfig   = "failed to init config: %v"
	FailedToConnectDB    = "failed to connect db: %v"
	FailedShutDownServer = "error occured on server shutting down: %v"
)
