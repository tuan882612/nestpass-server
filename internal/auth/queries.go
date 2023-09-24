package auth

const (
	UserCredsQuery string = "SELECT user_id, password FROM users WHERE email = $1"
	AddUserQuery   string = `
		INSERT INTO users 
			(user_id, email, name, password, registered, user_status) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	UpdateUserStatusQuery string = `
		UPDATE users
		SET user_status = 'active'
		WHERE user_id = $1`
)
