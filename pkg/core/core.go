package core

type Password struct {
	PasswordID  int64   `db:"password_id"`
	UserID      int64   `db:"user_id"`
	EncService  string  `db:"enc_service"`
	EncPassword string  `db:"enc_password"`
	EncLogin    *string `db:"enc_login"`
}

type User struct {
	UserID       int64  `db:"user_id"`
	Email        string `db:"email"`
	HashPassword string `db:"hash_password"`
	RegDate      string `db:"reg_date"`
}
