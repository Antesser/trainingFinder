package telegram

const (
	StateIdle          = "idle"
	StateAwaitingLogin = "awaiting_login"
	StateAwaitingPass  = "awaiting_password"
	StateAwaitingRole  = "awaiting_role"
)

type Session struct {
	State    string `json:"state"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
