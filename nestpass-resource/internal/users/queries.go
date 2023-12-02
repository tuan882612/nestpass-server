package users

const (
	GetUserQuery = `
		SELECT
			user_id,
			email,
			name,
			registered,
			user_status
		FROM users WHERE user_id = $1`
)
