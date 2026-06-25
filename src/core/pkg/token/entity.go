package token

type PayloadRole struct {
	ID   string
	Name string
	Type string
}

type Payload struct {
	UserID         string
	Role           PayloadRole
	PrivilageAdmin string
}

type Response struct {
	UserID    string
	RoleID    string
	Role      string
	RoleType  string
	Privilage string
}
