package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/dominikus1993/integrationtestcontainers-go/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestGetIdStage(t *testing.T) {
	article := model.NewArticle("testArticle", "http://test.com")
	articles := []model.Article{article}
	stage := getArticlesIds(articles)
	assert.NotEmpty(t, stage)
	assert.Equal(t, []model.ArticleId{"52c88c25-c646-5639-b444-a358092cc962"}, stage)
}

func TestIsNew(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// Arrange
	ctx := context.Background()
	mongoC, err := mongodb.StartContainer(ctx, mongodb.NewMongoContainerConfigurationBuilder().Build())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := mongoC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	connectionString, err := mongoC.ConnectionString(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("can't download mongo conectionstring, %w", err))
	}
	client, err := NewClient(ctx, connectionString, "Articles")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close(ctx)
	repo := NewMongoArticlesRepository(client)

	t.Run("Article not exists", func(t *testing.T) {
		// Act
		article := model.NewArticle("xd", "xDDDDD")
		stream := channels.FromSlice([]model.Article{article})
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
		stream := channels.FromSlice(articles)
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
		stream := channels.FromSlice(articles)
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
	mongoC, err := mongodb.StartContainer(ctx, mongodb.NewMongoContainerConfigurationBuilder().Build())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := mongoC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	connectionString, err := mongoC.ConnectionString(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("can't download mongo conectionstring, %w", err))
	}
	client, err := NewClient(ctx, connectionString, "Articles")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		client.Close(ctx)
	})
	repo := NewMongoArticlesRepository(client)

	t.Run("Article not exists", func(t *testing.T) {
		// Act
		article := model.NewArticle("xd", "xDDDDD")

		subject, err := repo.checkArticleExistence(context.TODO(), article)
		assert.NoError(t, err)
		assert.NotNil(t, subject)
		assert.Len(t, subject, 1)
		assert.False(t, subject[article.GetID()].existsInDb)
	})

	t.Run("Article exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle", "http://test.com")
		articles := []model.Article{article}
		err := repo.Save(ctx, articles)

		assert.Nil(t, err)

		subject, err := repo.checkArticleExistence(context.TODO(), articles...)
		assert.NoError(t, err)

		// Test
		assert.NotEmpty(t, subject)
		assert.Len(t, subject, 1)
		assert.True(t, subject[article.GetID()].existsInDb)
	})

	t.Run("Somearticle exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle2222", "http://test.com222222")
		articles := []model.Article{article}
		err := repo.Save(ctx, articles)
		newArticle := model.NewArticle("testArticle2", "http://test.com222222")
		articles = append(articles, newArticle)

		assert.NoError(t, err)
		subject, err := repo.checkArticleExistence(context.TODO(), articles...)
		// Test
		assert.NoError(t, err)
		assert.NotEmpty(t, subject)
		assert.Len(t, subject, 2)
		assert.False(t, subject[newArticle.GetID()].existsInDb)
		assert.True(t, subject[article.GetID()].existsInDb)
	})
}

func TestFilterOldArticles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// Arrange
	ctx := context.Background()
	mongoC, err := mongodb.StartContainer(ctx, mongodb.NewMongoContainerConfigurationBuilder().Build())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := mongoC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	connectionString, err := mongoC.ConnectionString(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("can't download mongo conectionstring, %w", err))
	}
	client, err := NewClient(ctx, connectionString, "Articles")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Close(ctx)
	})
	repo := NewMongoArticlesRepository(client)

	t.Run("Article not exists", func(t *testing.T) {
		// Act
		article := model.NewArticle("xd", "xDDDDD")

		subject, err := repo.filterOldArticles(context.TODO(), []model.Article{article})
		assert.NoError(t, err)
		assert.NotNil(t, subject)
		assert.Len(t, subject, 1)
	})

	t.Run("Article exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle", "http://test.com")
		articles := []model.Article{article}
		err := repo.Save(ctx, articles)

		assert.Nil(t, err)

		subject, err := repo.filterOldArticles(context.TODO(), articles)
		assert.NoError(t, err)

		// Test
		assert.Empty(t, subject)
	})

	t.Run("Somearticle exists in database", func(t *testing.T) {
		// Act
		article := model.NewArticle("testArticle2222", "http://test.com222222")
		articles := []model.Article{article}
		err := repo.Save(ctx, articles)
		newArticle := model.NewArticle("testArticle2", "http://test.com222222")
		articles = append(articles, newArticle)

		assert.NoError(t, err)
		subject, err := repo.filterOldArticles(context.TODO(), articles)
		// Test
		assert.NoError(t, err)
		assert.NotEmpty(t, subject)
		assert.Len(t, subject, 1)
		assert.Contains(t, subject, newArticle)
	})
}
