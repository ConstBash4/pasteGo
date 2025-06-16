package types

const (
	OperationSuccess    = 0
	OperationSuccessExp = "Success"

	ErrWrongCredentials    = 1001
	ErrWrongCredentialsExp = "Wrong username or password"

	ErrExistUser    = 1002
	ErrExistUserExp = "This user already exist"

	ErrUserNotFound    = 1003
	ErrUserNotFoundExp = "User not found"

	ErrUserSameCredentials    = 1004
	ErrUserSameCredentialsExp = "The same data when updating"

	ErrUserEmptyCredentials    = 1005
	ErrUserEmptyCredentialsExp = "Empty data for updating"

	ErrJWTProcessing    = 1101
	ErrJWTProcessingExp = "JWT processing error"

	ErrJWTExpired    = 1102
	ErrJWTExpiredExp = "JWT expired"

	ErrJWTNotFound    = 1103
	ErrJWTNotFoundExp = "JWT not found"

	ErrGetCookies    = 1201
	ErrGetCookiesExp = "Cookie reading error"

	ErrEmptyPaste    = 2001
	ErrEmptyPasteExp = "Paste cannot be empty"

	ErrPasteNotFound    = 2002
	ErrPasteNotFoundExp = "Paste not found"

	ErrNotPublicPaste    = 2003
	ErrNotPublicPasteExp = "Access only for authorized users"

	ErrPasswordPaste    = 2004
	ErrPasswordPasteExp = "This paste requires a password"

	ErrWrongPasswordPaste    = 2005
	ErrWrongPasswordPasteExp = "Wrong password"

	ErrServer    = 5000
	ErrServerExp = "Server problem"
)
