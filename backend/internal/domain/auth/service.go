package auth

import "github.com/euandresimoes/ecom-go/backend/internal/models"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(data models.UserRegisterModel) error {
	return s.repo.Register(data)
}

func (s *Service) Login(data models.UserLoginModel) (string, error) {
	return s.repo.Login(data)
}

func (s *Service) Profile(id float64) (models.UserPublicModel, error) {
	return s.repo.Profile(id)
}
