package repository

import (
	"database/sql"
	"errors"

	"dryka.pl/trader/internal/domain/trade/model"
	"dryka.pl/trader/internal/domain/trade/repository"
	_ "github.com/mattn/go-sqlite3"
)

var ErrPersistencePrepareError = errors.New("persistence prepare error")
var ErrPersistenceCannotAddLog = errors.New("persistence cannot add log")

type sqliteAuditRepository struct {
	db *sql.DB
}

func (s sqliteAuditRepository) Log(log model.AuditLog) error {
	stmt, err := s.db.Prepare(
		`INSERT INTO audit_log
   	(source_quantity, source_currency, target_currency, target_quantity, timestamp, update_id)
   	VALUES (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return ErrPersistencePrepareError
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	_, err = stmt.Exec(log.SourceQuantity, log.SourceCurrency, log.TargetCurrency, log.TargetQuantity, log.Timestamp, log.UpdateId)
	if err != nil {
		return ErrPersistenceCannotAddLog
	}

	return nil
}

func NewAuditRepository(db *sql.DB) repository.Audit {
	return &sqliteAuditRepository{
		db: db,
	}
}
