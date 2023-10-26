package passwords

const (
	GetKDFDataQuery = `
	SELECT
		password,
		salt
	FROM users WHERE user_id = $1`

	GetAllPasswordsNonPagedQuery = `
	SELECT * FROM passwords
	WHERE user_id = $1`

	GetAllPasswordsQuery = `
	SELECT * FROM passwords
	WHERE user_id = $1 AND password_id > $2
	ORDER BY password_id ASC
	LIMIT $3`

	GetAllPasswordsByCategoryQuery = `
	SELECT * FROM passwords
	WHERE user_id = $1 AND category_id = $2 AND password_id > $3
	ORDER BY password_id ASC
	LIMIT $4`

	GetPasswordQuery = `
	SELECT * FROM passwords
	WHERE user_id = $1 AND password_id = $2 AND category_id = $3`

	CreatePasswordQuery = `
	INSERT INTO passwords (
		password_id, user_id, category_id, website, nonce, encrypted
	) VALUES ($1, $2, $3, $4, $5, $6)`

	UpdatePasswordQuery = `
	UPDATE passwords SET website = $1, nonce = $2, encrypted = $3
	WHERE password_id = $4 AND category_id = $5 AND user_id = $6`

	DeletePasswordQuery = `
	DELETE FROM passwords
	WHERE password_id = $1 AND category_id = $2 AND user_id = $3`
)
