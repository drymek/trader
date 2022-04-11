package repository

type AccountRepository interface {
	Create(account interface{}) error
	Find(s string) (interface{}, error)
	Delete(id string) error
	Update(entity interface{}) error
}
