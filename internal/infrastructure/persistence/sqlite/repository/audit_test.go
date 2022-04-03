package repository

import (
	"database/sql"
	"testing"
	"time"

	"dryka.pl/trader/internal/domain/trade/model"
	"dryka.pl/trader/internal/infrastructure/persistence/sqlite"
	"dryka.pl/trader/tests/database"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB *sql.DB
}

func TestSuite(t *testing.T) {
	s := new(Suite)

	err := database.CreateFromTemplate(
		"../../../../../database/sqlite/database.sqlite.template",
		"../../../../../database/sqlite/database_test.sqlite",
	)
	if err != nil {
		t.Fatal(err)
	}

	connection, err := sqlite.NewConnection("../../../../../database/sqlite/database_test.sqlite")
	if err != nil {
		t.Fatal(err)
	}

	s.DB = connection
	suite.Run(t, s)
}

func (s *Suite) TestValidLog() {
	var count int

	err := NewAuditRepository(s.DB).Log(model.AuditLog{
		SourceQuantity: "42",
		SourceCurrency: "BNB",
		TargetCurrency: "USDT",
		TargetQuantity: "5",
		Timestamp:      time.Now(),
		UpdateId:       123,
	})
	s.NoError(err)

	exec := s.DB.QueryRow("SELECT count(*) FROM audit_log")
	err = exec.Scan(&count)
	s.NoError(err)
	s.Equal(1, count)
}
