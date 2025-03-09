package postgres

const (
	queryRegistration = "INSERT INTO users (email, hash_password, reg_date) VALUES (:email, :hash_password, :reg_date) RETURNING user_id"
	queryLogIn        = "SELECT user_id, email, hash_password FROM users WHERE email = $1"
	queryDelUser      = "DELETE FROM users WHERE user_id = $1"
	queryUpdateUser   = "UPDATE users SET email = :email, hash_password = :hash_password WHERE user_id = :user_id"
	queryGetUserByID  = "SELECT user_id, email, hash_password FROM users WHERE user_id = $1"

	queryInserNewPassword    = "INSERT INTO passwords (user_id, enc_service, enc_password) VALUES (:user_id, :enc_service, :enc_password)"
	querySelectUserPasswords = "SELECT enc_service, enc_password FROM passwords WHERE user_id = $1"
	queryDelPassword         = "DELETE FROM passwords WHERE password_id = $1"
	queryGetPasswordByID     = "SELECT password_id, enc_service, enc_password FROM passwords WHERE password_id = $1"
)
