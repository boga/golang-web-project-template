package model

type User struct {
	ID             int
	Name           *string
	AuthIdentities []AuthIdentity
}

type AuthIdentity struct {
	ID       int
	Password *string
	UID      string
	User     User
	UserID   int `db:"user_id"`
}
