package repository

import (
	"database/sql"

	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/repository"
)

type sqliteAccountRepository struct {
	db *sql.DB
}

func (a *sqliteAccountRepository) Update(entity interface{}) error {
	account := entity.(*model.Account)

	stmt, err := a.db.Prepare(
		`UPDATE account SET owner = ? , balance = ?, currency = ?, account_number = ? WHERE custom_id = ?`,
	)
	if err != nil {
		return ErrPersistencePrepareError
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)
	_, err = stmt.Exec(account.Owner, account.Balance, account.Currency, account.AccountNumber, account.ID)
	if err != nil {
		return ErrPersistenceError
	}

	return nil
}

func (a *sqliteAccountRepository) Delete(id string) error {
	stmt, err := a.db.Prepare(
		`DELETE FROM account WHERE custom_id = ?`,
	)
	if err != nil {
		return ErrPersistencePrepareError
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	_, err = stmt.Exec(id)
	if err != nil {
		return ErrPersistenceError
	}

	return nil
}

func (a *sqliteAccountRepository) Find(id string) (interface{}, error) {
	stmt, err := a.db.Prepare(
		`SELECT custom_id, owner, balance, currency, account_number FROM account WHERE custom_id = ?`,
	)
	if err != nil {
		return nil, ErrPersistencePrepareError
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)
	account := &model.Account{}
	row := stmt.QueryRow(id)
	if err := row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.AccountNumber); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrAccountNotFound
		}
		return nil, ErrPersistenceError
	}

	return account, nil
}

func (a *sqliteAccountRepository) Create(m interface{}) error {
	account := m.(*model.Account)
	stmt, err := a.db.Prepare(
		`INSERT INTO account
   	(custom_id, owner, balance, currency, account_number)
   	VALUES (?, ?, ?, ?, ?)`)

	if err != nil {
		return ErrPersistencePrepareError
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	_, err = stmt.Exec(account.ID, account.Owner, account.Balance, account.Currency, account.AccountNumber)

	if err != nil {
		return ErrPersistenceCannotAddAccount
	}

	return nil
}

func NewAccountRepository(db *sql.DB) repository.AccountRepository {
	return &sqliteAccountRepository{
		db: db,
	}
}
