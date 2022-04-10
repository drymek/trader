package repository

import (
	"context"

	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountRepository struct {
	collection *mongo.Collection
}

func (a *accountRepository) Update(entity interface{}) error {
	id := entity.(*model.Account).ID
	filter := bson.D{{"id", id}}

	_, err := a.collection.ReplaceOne(context.TODO(), filter, entity)
	return err
}

func (a *accountRepository) Delete(id string) error {
	filter := bson.D{{"id", id}}

	_, err := a.collection.DeleteOne(context.TODO(), filter)
	return err
}

func (a *accountRepository) Find(ID string) (interface{}, error) {
	filter := bson.D{{"id", ID}}

	var account model.Account
	err := a.collection.FindOne(context.Background(), filter).Decode(&account)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrAccountNotFound
		}
		return nil, err
	}
	return &account, err
}

func (a *accountRepository) Create(model interface{}) error {
	_, err := a.collection.InsertOne(context.TODO(), model)
	return err
}

func NewAccountRepository(client *mongo.Client) repository.AccountRepository {
	return &accountRepository{
		collection: client.Database("myDB").Collection("accounts"),
	}
}
