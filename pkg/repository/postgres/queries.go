package postgres

const (
	queryRegistration = "INSERT INTO users (email, hash_password, reg_date) VALUES (:email, :hash_password, :reg_date) RETURNING user_id"
	queryLogIn        = "SELECT user_id, email, hash_password FROM users WHERE email = $1"
)
