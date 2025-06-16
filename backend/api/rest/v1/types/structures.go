package types

const (
	CookieAccessToken  = "access_token"
	CookieRefreshToken = "refresh_token"
	CookieUsername     = "username"
	CookieExp          = "exp"
)

var SecretKey []byte

type APIResponse struct {
	Code        int    `json:"code"`
	Explanation string `json:"explanation"`
	Message     any    `json:"message,omitempty"`
}

type User struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Paste struct {
	Id          string `json:"id"`
	Author      string `json:"author"`
	Created     int64  `json:"created,omitempty"`
	Updated     int64  `json:"updated,omitempty"`
	ExpTime     int64  `json:"expTime"`
	Lifetime    string `json:"lifetime,omitempty"`
	Text        string `json:"text"`
	Password    string `json:"password"`
	Public      bool   `json:"public"`
	HasPassword bool   `json:"hasPassword"`
}

type PasteList struct {
	Pastes []Paste `json:"pastes"`
}

type PastePassword struct {
	Password string `json:"password,omitempty"`
}

type Tokens struct {
	RefreshToken string
	AccessToken  string
}
