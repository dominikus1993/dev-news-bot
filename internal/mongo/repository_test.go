package mongo

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/integrationtestcontainers-go"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSave(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		repo := NewMongoArticlesRepository(&client)
		article := model.NewArticle("testArticle", "http://testarticle.com")
		err := repo.Save(context.Background(), []model.Article{article})
		assert.Nil(mt, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "duplicate key error",
		}))
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		repo := NewMongoArticlesRepository(&client)

		article := model.NewArticle("testArticle", "http://test.com")
		err := repo.Save(context.Background(), []model.Article{article})
		assert.NotNil(mt, err)
	})
}

func TestIsNew(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// Arrange
	ctx := context.Background()
	mongoC, err := integrationtestcontainers.StartMongoDbContainer(ctx, integrationtestcontainers.DefaultMongoContainerConfiguration)
	if err != nil {
		t.Fatal(err)
	}
	defer mongoC.Terminate(ctx)
	client, err := NewClient(ctx, mongoC.ConnectionString, "Articles")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close(ctx)
	repo := NewMongoArticlesRepository(client)

	t.Run("Article not exists", func(t *testing.T) {
		// Act
		article := model.NewArticle("xd", "xDDDDD")
		stream := channels.FromSlice(article)
		res := repo.FilterNew(context.TODO(), stream)
		subject := channels.ToSlice(res)
		assert.NotNil(t, subject)
		assert.NotEmpty(t, subject)
		assert.Len(t, subject, 1)
	})

	t.Run("Article exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle", "http://test.com")
		articles := []model.Article{article}
		err := repo.Save(ctx, articles)

		assert.Nil(t, err)
		stream := channels.FromSlice(articles...)
		res := repo.FilterNew(context.TODO(), stream)
		subject := channels.ToSlice(res)
		// Test
		assert.Empty(t, subject)
	})

	t.Run("Somearticle exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle2222", "http://test.com22")
		articles := []model.Article{article}
		err := repo.Save(ctx, articles)
		newArticle := model.NewArticle("testArticle2", "http://test.com2")
		articles = append(articles, newArticle)

		assert.Nil(t, err)
		stream := channels.FromSlice(articles...)
		res := repo.FilterNew(context.TODO(), stream)
		subject := channels.ToSlice(res)
		// Test
		assert.NotEmpty(t, subject)
		assert.Len(t, subject, 1)
		element := subject[0]
		assert.Equal(t, newArticle.GetID(), element.GetID())
	})
}

func TestGetIdsThatExistsInDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// Arrange
	ctx := context.Background()
	mongoC, err := integrationtestcontainers.StartMongoDbContainer(ctx, integrationtestcontainers.DefaultMongoContainerConfiguration)
	if err != nil {
		t.Fatal(err)
	}
	defer mongoC.Terminate(ctx)
	client, err := NewClient(ctx, mongoC.ConnectionString, "Articles")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close(ctx)
	repo := NewMongoArticlesRepository(client)

	t.Run("Article not exists", func(t *testing.T) {
		// Act
		article := model.NewArticle("xd", "xDDDDD")

		subject, err := repo.getIdsThatExistsInDatabase(context.TODO(), article)
		assert.NoError(t, err)
		assert.NotNil(t, subject)
		assert.Empty(t, subject)
	})

	t.Run("Article exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle", "http://test.com")
		articles := []model.Article{article}
		err := repo.Save(ctx, articles)

		assert.Nil(t, err)

		subject, err := repo.getIdsThatExistsInDatabase(context.TODO(), articles...)
		assert.NoError(t, err)

		// Test
		assert.NotEmpty(t, subject)
		assert.Len(t, subject, 1)
	})

	t.Run("Somearticle exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle2222", "http://test.com22")
		articles := []model.Article{article}
		err := repo.Save(ctx, articles)
		newArticle := model.NewArticle("testArticle2", "http://test.com2")
		articles = append(articles, newArticle)

		assert.NoError(t, err)
		subject, err := repo.getIdsThatExistsInDatabase(context.TODO(), articles...)
		// Test
		assert.NoError(t, err)
		assert.NotEmpty(t, subject)
		assert.Len(t, subject, 1)
		element := subject[0]
		assert.Equal(t, newArticle.GetID(), element)
	})
}
