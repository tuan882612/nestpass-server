package auth

const (
	UserCreds string = "SELECT id, password FROM users WHERE email = $1"
	AddUser   string = `
		INSERT INTO users 
			(user_id, email, name, password, registered, user_status) 
		VALUES ($1, $2, $3, $4, $5, $6)`
)
