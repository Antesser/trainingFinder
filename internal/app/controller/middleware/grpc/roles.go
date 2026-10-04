package grpc

const (
	User  string = "user"
	Coach string = "coach"
	Admin string = "admin"
)

var (
	roles = map[string]string{
		"5": Admin,
		"6": Coach,
		"7": User,
	}
)

func IsCoach(input string) bool {
	return roles[input] == Coach
}

func IsUser(input string) bool {
	return roles[input] == User
}
