package exception

type ErrorCode uint8

const (
	UNKNOWN ErrorCode = 255
)

const (
	// auth
	NO_TOKEN       ErrorCode = 0
	INVALID_TOKEN  ErrorCode = 1
	INCORRECT_AUTH ErrorCode = 2

	// file
	INVALID_FOLDER_ID ErrorCode = 11

	// user
	EMAIL_EXISTS ErrorCode = 21
)
