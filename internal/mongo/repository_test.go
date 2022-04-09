package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSave(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		repo := NewMongoArticlesRepository(&client)
		article := model.NewArticle("testArticle", "http://test.com")
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

type mongoContainer struct {
	testcontainers.Container
	ConnectionString string
}

func setupMongo(ctx context.Context) (*mongoContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := mongoC.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		return nil, err
	}
	mongoConnection := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	return &mongoContainer{mongoC, mongoConnection}, nil
}

func TestIsNew(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// Arrange
	ctx := context.Background()
	mongoC, err := setupMongo(ctx)
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
		isNew, err := repo.IsNew(ctx, &article)

		// Test
		assert.Nil(t, err)
		assert.True(t, isNew)
	})

	t.Run("Article exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle", "http://test.com")
		err := repo.Save(ctx, []model.Article{article})

		assert.Nil(t, err)

		isNew, err := repo.IsNew(ctx, &article)

		// Test
		assert.Nil(t, err)
		assert.False(t, isNew)
	})
}
