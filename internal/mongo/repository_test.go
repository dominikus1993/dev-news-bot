package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSave(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		repo := NewMongoArticlesRepository(&client)

		article := model.Article{Title: "test"}
		err := repo.Save(mtest.Background, []model.Article{article})
		assert.Nil(mt, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "duplicate key error",
		}))
		client := MongoClient{mt.Client, mt.DB, mt.Coll}
		repo := NewMongoArticlesRepository(&client)

		article := model.Article{Title: "test"}
		err := repo.Save(mtest.Background, []model.Article{article})
		assert.NotNil(mt, err)
	})
}

func TestIsNew(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017"),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	defer mongoC.Terminate(ctx)
	host, err := mongoC.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		t.Error(err)
	}
	mongoConnection := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	client, err := NewClient(ctx, mongoConnection, "Articles")
	if err != nil {
		t.Error(err)
	}
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

		isNew, err := repo.IsNew(ctx, &article)

		// Test
		assert.Nil(t, err)
		assert.False(t, isNew)
	})
}
