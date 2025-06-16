package typesDB

type UserRecord struct {
	Id       string //UUID
	Username string
	Password string
}

type PasteRecord struct {
	Id       string //UUID
	UserId   string
	Text     string
	Created  int64
	Updated  int64
	Lifetime int64
	Password string
	Public   int
}

type TokenRecord struct {
	RefreshToken string
	UserId       string
}

const (
	UsersTable  = "users"
	PastesTable = "pastes"
)

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func IntToBool(i int) bool {
	return i != 0
}
