package redis

type Session struct {
	State    string `json:"state"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

const (
	KeyPrefix = "bot:session:"
)
