package service

import (
	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/repository"
)

type CrudService interface {
	Create(entity interface{}) error
	Fetch(id string) (interface{}, error)
	UpdateOrCreate(account interface{}) (bool, error)
	Delete(id string) error
}

type AccountService struct {
	repository repository.AccountRepository
}

func (s *AccountService) Delete(id string) error {
	return s.repository.Delete(id)
}

func (s *AccountService) UpdateOrCreate(entity interface{}) (bool, error) {
	id := entity.(model.Entity).GetID()
	_, err := s.repository.Find(id)

	if err == repository.ErrAccountNotFound {
		return false, s.repository.Create(entity)
	}

	if err != nil {
		return false, err
	}

	return true, s.repository.Update(entity)
}

func (s *AccountService) Fetch(id string) (interface{}, error) {
	return s.repository.Find(id)
}

func (s *AccountService) Create(entity interface{}) error {
	return s.repository.Create(entity)
}

func NewAccountService(repository repository.AccountRepository) CrudService {
	return &AccountService{
		repository: repository,
	}
}
