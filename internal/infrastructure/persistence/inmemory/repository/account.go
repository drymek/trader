package repository

import (
	"sync"

	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/repository"
)

type AccountRepository struct {
	mu sync.Mutex
	m  map[string]*model.Account
}

func (a *AccountRepository) Update(entity interface{}) error {
	_, err := a.Find(entity.(model.Entity).GetID())
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.m[entity.(model.Entity).GetID()] = entity.(*model.Account)

	return nil
}

func (a *AccountRepository) Delete(id string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.m[id]; !ok {
		return repository.ErrAccountNotFound
	}

	delete(a.m, id)
	return nil
}

func (a *AccountRepository) Find(ID string) (interface{}, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if account, ok := a.m[ID]; ok {
		return account, nil
	}

	return nil, repository.ErrAccountNotFound
}

func (a *AccountRepository) Create(m interface{}) error {
	account, ok := m.(*model.Account)
	if !ok {
		return repository.ErrInvalidModel
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.m[account.ID]; ok {
		return repository.ErrAccountAlreadyExists
	}

	a.m[account.ID] = account

	return nil
}

func NewAccountRepository() repository.AccountRepository {
	return &AccountRepository{
		mu: sync.Mutex{},
		m:  make(map[string]*model.Account),
	}
}
