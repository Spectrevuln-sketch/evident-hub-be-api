package context

type ctxKey string

const (
	UserIDKey ctxKey = "userID"
	RoleIDKey ctxKey = "roleID"
	RoleKey   ctxKey = "role"
)

func (s ctxKey) String() string {
	return string(s)
}
