package roles

const (
	User  string = "user"
	Coach string = "coach"
	Admin string = "admin"
)

var (
	Roles = map[string]string{
		"5": Admin,
		"6": Coach,
		"7": User,
	}
)

func RoleKey(role string) (string, bool) {
	for k, v := range Roles {
		if v == role {

			return k, true
		}
	}

	return "", false
}

func IsCoach(input string) bool {
	return Roles[input] == Coach
}

func IsUser(input string) bool {
	return Roles[input] == User
}
