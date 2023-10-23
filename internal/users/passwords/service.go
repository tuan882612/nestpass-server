package passwords

type service struct {
	repo *repository
}

func NewService(repo *repository) *service {
	return &service{
		repo: repo,
	}
}
