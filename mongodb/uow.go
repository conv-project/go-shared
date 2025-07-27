package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUnitOfWork struct {
	client *mongo.Client
}

func NewMongoUnitOfWork(client *mongo.Client) *MongoUnitOfWork {
	return &MongoUnitOfWork{client: client}
}

func (u *MongoUnitOfWork) Start(ctx context.Context) (mongo.Session, error) {
	session, err := u.client.StartSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (u *MongoUnitOfWork) Commit(ctx context.Context, session mongo.Session) error {
	return session.CommitTransaction(ctx)
}

func (u *MongoUnitOfWork) Abort(ctx context.Context, session mongo.Session) error {
	return session.AbortTransaction(ctx)
}
