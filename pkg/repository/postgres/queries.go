package postgres

const (
	queryRegistration = "INSERT INTO users (email, hash_password, reg_date) VALUES (:email, :hash_password, :reg_date) RETURNING user_id"
	queryLogIn        = "SELECT user_id, email, hash_password FROM users WHERE email = $1"
	queryDelUser      = "DELETE FROM users WHERE user_id = $1"
	queryUpdateUser   = "UPDATE users SET email = :email, hash_password = :hash_password WHERE user_id = :user_id"
)
