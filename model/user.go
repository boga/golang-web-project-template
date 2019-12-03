package model

type User struct {
	ID              int     `db:"id" json:"id"`
	Name            *string `db:"name" json:"name"`
	AuthIdentities  []AuthIdentity
	TOTPBackupCodes *string `db:"totp_backup_codes" json:"totp_backup_codes"`
	TOTPEnabled     bool    `db:"totp_enabled" json:"totp_enabled"`
	TOTPSecret      *string `db:"totp_secret" json:"totp_secret"`
}

type AuthIdentity struct {
	ID       int
	Password *string
	UID      string
	User     User
	UserID   int `db:"user_id"`
}
