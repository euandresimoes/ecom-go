package product

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(data *CategoryCreateDto) (CategoryModel, error) {
	return s.repo.CreateCategory(data)
}

func (s *Service) GetAllCategories() ([]CategoryModel, error) {
	return s.repo.GetAllCategories()
}

func (s *Service) DeleteCategory(id int) (CategoryModel, error) {
	return s.repo.DeleteCategory(id)
}

func (s *Service) Create(data *ProductCreateDto) (ProductModel, error) {
	return s.repo.Create(data)
}

func (s *Service) Delete(id int) (ProductModel, error) {
	return s.repo.Delete(id)
}

func (s *Service) Update(id int, data *ProductUpdateDto) (ProductModel, error) {
	return s.repo.Update(id, data)
}

func (s *Service) GetAll() ([]ProductModel, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (ProductModel, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetByPublicID(publicID string) (ProductModel, error) {
	return s.repo.GetByPublicID(publicID)
}
