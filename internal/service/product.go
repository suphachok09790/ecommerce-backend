package service

import (
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"errors"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]model.Product, error) {
	return s.repo.FindAll()
}

func (s *ProductService) GetByID(id uint) (model.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Create(p *model.Product) error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Price < 0 {
		return errors.New("price cannot be negative")
	}
	if p.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	return s.repo.Create(p)
}

func (s *ProductService) Update(id uint, input *model.Product) (model.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return product, errors.New("product not found")
	}

	if input.Name != "" {
		product.Name = input.Name
	}
	if input.Description != "" {
		product.Description = input.Description
	}
	if input.Price >0 {
		product.Price = input.Price
	}
	if input.Stock >= 0 {
		product.Stock = input.Stock
	}

	err = s.repo.Update(&product)
	return product, err
}

func (s *ProductService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}
	return s.repo.Delete(id)
}