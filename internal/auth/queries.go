package auth

const (
	UserCredsQuery string = `
		SELECT 
			user_id, password, user_status 
		FROM users 
		WHERE email = $1`
	GetUserPasswordQuery string = `
		SELECT password
		FROM users
		WHERE user_id = $1`
	AddUserQuery   string = `
		INSERT INTO users 
			(user_id, email, name, password, registered, user_status, salt) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	UpdateUserStatusQuery string = `
		UPDATE users
		SET user_status = 'active'
		WHERE user_id = $1`
	UpdateUserPasswordQuery string = `
		UPDATE users
		SET password = $2
		WHERE user_id = $1`
)
