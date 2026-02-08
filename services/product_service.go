package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) Create(input *models.ProductInput) (*models.Product, error) {
	return s.repo.Create(input)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(id int, input *models.ProductInput) (*models.Product, error) {
	return s.repo.Update(id, input)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
