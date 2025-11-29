package auth

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) register(data UserRegisterModel) error {
	return s.repo.Register(data)
}

func (s *Service) login(data UserLoginModel) (string, error) {
	return s.repo.Login(data)
}
