package auth

type roleStruct struct {
	ID   int    `db:"id"`
	Role string `db:"role"`
}
