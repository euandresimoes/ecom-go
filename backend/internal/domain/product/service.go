package product

import "github.com/euandresimoes/ecom-go/backend/internal/models"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(data *models.CategoryCreateDto) (models.CategoryModel, error) {
	return s.repo.CreateCategory(data)
}

func (s *Service) GetAllCategories() ([]models.CategoryModel, error) {
	return s.repo.GetAllCategories()
}

func (s *Service) DeleteCategory(id int) (models.CategoryModel, error) {
	return s.repo.DeleteCategory(id)
}

func (s *Service) Create(data *models.ProductCreateDto) (models.ProductModel, error) {
	return s.repo.Create(data)
}

func (s *Service) Delete(id int) (models.ProductModel, error) {
	return s.repo.Delete(id)
}

func (s *Service) Update(id int, data *models.ProductUpdateDto) (models.ProductModel, error) {
	return s.repo.Update(id, data)
}

func (s *Service) GetAll() ([]models.ProductModel, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (models.ProductModel, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetByPublicID(publicID string) (models.ProductModel, error) {
	return s.repo.GetByPublicID(publicID)
}
