package postgres

const (
	queryRegistration = "INSERT INTO users (email, hash_password, reg_date) VALUES (:email, :hash_password, :reg_date) RETURNING user_id"
)
