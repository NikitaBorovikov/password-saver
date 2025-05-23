package postgres

const (
	//user table queries
	queryRegistration   = "INSERT INTO users (email, hash_password, reg_date) VALUES (:email, :hash_password, :reg_date) RETURNING user_id"
	queryLogIn          = "SELECT user_id, email, hash_password FROM users WHERE email = $1"
	queryDelUser        = "DELETE FROM users WHERE user_id = $1"
	queryUpdateUser     = "UPDATE users SET email = :email, hash_password = :hash_password WHERE user_id = :user_id"
	querySelectUserByID = "SELECT user_id, email, hash_password FROM users WHERE user_id = $1"

	//password table queries
	queryInserNewPassword    = "INSERT INTO passwords (user_id, enc_service, enc_password, enc_login) VALUES (:user_id, :enc_service, :enc_password, :enc_login)"
	querySelectUserPasswords = "SELECT password_id, enc_service, enc_password, enc_login FROM passwords WHERE user_id = $1"
	queryDelPassword         = "DELETE FROM passwords WHERE password_id = $1 AND user_id = $2"
	querySelectPasswordByID  = "SELECT password_id, enc_service, enc_password, enc_login FROM passwords WHERE password_id = $1 AND user_id = $2"
	queryUpdatePassword      = "UPDATE passwords SET enc_service = :enc_service, enc_password = :enc_password, enc_login = :enc_login WHERE password_id = :password_id AND user_id = :user_id"
)
