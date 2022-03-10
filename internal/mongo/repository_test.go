package mongo

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSaveWhenSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		repo := NewMongoArticlesRepository(&client)

		article := model.Article{Title: "test"}
		err := repo.Save(mtest.Background, []model.Article{article})
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
		repo := NewMongoArticlesRepository(&client)

		article := model.Article{Title: "test"}
		err := repo.Save(mtest.Background, []model.Article{article})
		assert.NotNil(t, err)
	})
}

func TestIsNewWhenArticleExistsInDb(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		repo := NewMongoArticlesRepository(&client)
		article := model.Article{Title: "test"}
		err := repo.Save(mtest.Background, []model.Article{article})
		assert.Nil(t, err)
		isNew, err := repo.IsNew(mtest.Background, &article)
		assert.Nil(t, err)
		assert.False(t, isNew)
	})
}

func TestIsNewWhenArticleNotExistsInDb(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		ctx := context.TODO()
		repo := NewMongoArticlesRepository(&client)
		article := model.Article{Title: "test"}
		err := repo.Save(ctx, []model.Article{article})
		assert.Nil(t, err)
		article2 := model.Article{Title: "test2"}
		isNew, err := repo.IsNew(ctx, &article2)
		assert.Nil(t, err)
		assert.True(t, isNew)
	})
}
