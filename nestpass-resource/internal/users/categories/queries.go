package categories

const (
	GetAllCategoriesQuery = `
	SELECT * FROM categories
	WHERE user_id = $1 AND category_id > $2
	ORDER BY category_id ASC
	LIMIT $3`

	GetNameCategoryQuery = `
	SELECT * FROM categories
	WHERE name = $1 AND user_id = $2`

	GetUUIDCategoryQuery = `
	SELECT * FROM categories
	WHERE category_id = $1 AND user_id = $2`

	InsertCategoryQuery = `
	INSERT INTO categories (category_id, user_id, name, description)
	VALUES ($1, $2, $3, $4)`

	UpdateCategoryQuery = `
	UPDATE categories SET name = $1, description = $2
	WHERE category_id = $3 AND user_id = $4`

	DeleteCategoryQuery = `
	DELETE FROM categories
	WHERE category_id = $1 AND user_id = $2`
)
