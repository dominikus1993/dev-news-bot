package mongo

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func fakeClient(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
}

func TestSaveWhenSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		ctx := context.TODO()
		repo := NewMongoArticlesRepository(&client)

		article := model.Article{Title: "test"}
		err := repo.Save(ctx, []model.Article{article})
		assert.Nil(t, err)
	})
}

func TestSaveWhenError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "duplicate key error",
		}))
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		ctx := context.TODO()
		repo := NewMongoArticlesRepository(&client)

		article := model.Article{Title: "test"}
		err := repo.Save(ctx, []model.Article{article})
		assert.NotNil(t, err)
	})
}
