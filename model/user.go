package model

type User struct {
	ID             int
	Name           *string
	AuthIdentities []AuthIdentity
}

type AuthIdentity struct {
	ID   int
	UID  string
	User User
}
