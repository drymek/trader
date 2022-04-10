package repository

import (
	"context"
	"sync"
	"testing"
	"time"

	"dryka.pl/trader/internal/application/config"
	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/repository"
	inmemory "dryka.pl/trader/internal/infrastructure/persistence/inmemory/repository"
	mongox "dryka.pl/trader/internal/infrastructure/persistence/mongo/repository"
	sqlitex "dryka.pl/trader/internal/infrastructure/persistence/sqlite"
	sqlite "dryka.pl/trader/internal/infrastructure/persistence/sqlite/repository"
	"dryka.pl/trader/tests/database"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	repositories map[string]repository.AccountRepository
	mu           sync.Mutex
}

func TestServiceSuite(t *testing.T) {
	s := new(RepositoryIntegrationSuite)
	err := database.CreateFromTemplate(
		"../../../../../database/sqlite/database.sqlite.template",
		"../../../../../database/sqlite/database_test.sqlite",
	)
	if err != nil {
		t.Fatal(err)
	}
	connection, _ := sqlitex.NewConnection("../../../../../database/sqlite/database_test.sqlite")
	c, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(c.GetMongoURI()))
	if err != nil {
		t.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)
	s.repositories = make(map[string]repository.AccountRepository)
	s.repositories["mongo"] = mongox.NewAccountRepository(client)
	s.repositories["inmemory"] = inmemory.NewAccountRepository()
	s.repositories["sqlite"] = sqlite.NewAccountRepository(connection)
	_ = mongox.NewAccountRepository(client)
	_ = inmemory.NewAccountRepository()
	_ = sqlite.NewAccountRepository(connection)
	s.mu = sync.Mutex{}
	suite.Run(t, s)
}

func (s *RepositoryIntegrationSuite) SetupTest() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.repositories {
		account, err := s.repositories[i].Find("123")
		if err != nil {
			return
		}
		err = s.repositories[i].Delete(account.(*model.Account).ID)
		s.NoError(err)
	}
}

func (s *RepositoryIntegrationSuite) TestCreateAndFetch() {
	for i := range s.repositories {
		s.T().Run(i, func(t *testing.T) {
			s.mu.Lock()
			defer s.mu.Unlock()
			account := &model.Account{
				ID:            "123",
				Owner:         "Marcin Dryka",
				Balance:       "100.0",
				Currency:      "PLN",
				AccountNumber: 123456789,
			}

			err := s.repositories[i].Create(account)
			s.NoError(err)
			got, err := s.repositories[i].Find(account.ID)
			s.NoError(err)
			s.Equal(account, got)
		})
	}
}

func (s *RepositoryIntegrationSuite) TestDelete() {
	for i := range s.repositories {
		s.T().Run(i, func(t *testing.T) {
			s.mu.Lock()
			defer s.mu.Unlock()
			account := &model.Account{
				ID:            "123",
				Owner:         "Marcin Dryka",
				Balance:       "100.0",
				Currency:      "PLN",
				AccountNumber: 123456789,
			}

			err := s.repositories[i].Create(account)
			s.NoError(err)

			err = s.repositories[i].Delete(account.ID)
			s.NoError(err)

			_, err = s.repositories[i].Find(account.ID)
			s.Error(err)
			s.Equal(repository.ErrAccountNotFound, err)
		})
	}
}

func (s *RepositoryIntegrationSuite) TestUpdate() {
	for i := range s.repositories {
		s.T().Run(i, func(t *testing.T) {
			s.mu.Lock()
			defer s.mu.Unlock()
			account := &model.Account{
				ID:            "123",
				Owner:         "Marcin Dryka",
				Balance:       "100.0",
				Currency:      "PLN",
				AccountNumber: 123456789,
			}

			update := &model.Account{
				ID:            "123",
				Owner:         "John Doe",
				Balance:       "300.0",
				Currency:      "USD",
				AccountNumber: 123456789,
			}

			err := s.repositories[i].Create(account)
			s.NoError(err)

			err = s.repositories[i].Update(update)
			s.NoError(err)

			got, err := s.repositories[i].Find(account.ID)
			s.NoError(err)

			want := update
			if diff := cmp.Diff(want, got); diff != "" {
				s.Failf("value mismatch", "(-want +got):\n%v", diff)
			}
		})
	}
}
