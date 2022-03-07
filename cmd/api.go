package cmd

import (
	"context"
	"flag"

	"github.com/dominikus1993/dev-news-bot/internal/mongo"
	"github.com/dominikus1993/dev-news-bot/pkg/common"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
	"github.com/dominikus1993/dev-news-bot/pkg/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"
)

type RunDevNewsApi struct {
	mongoConnectionString string
}

func (*RunDevNewsApi) Name() string     { return "run-api" }
func (*RunDevNewsApi) Synopsis() string { return "run api" }
func (*RunDevNewsApi) Usage() string {
	return `go run . run-api"`
}

func (p *RunDevNewsApi) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.mongoConnectionString, "mongo-connection-string", "", "mongo connection string")
}

func (p *RunDevNewsApi) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	engine := html.New("./public", ".html")
	log.WithField("mongo", p.mongoConnectionString).Infoln("Start DevNews Api")
	mongodbClient, err := mongo.NewClient(ctx, p.mongoConnectionString)
	if err != nil {
		log.WithError(err).Error("can't create mongodb client")
		return subcommands.ExitFailure
	}
	defer mongodbClient.Close(ctx)
	repo := mongo.NewMongoArticlesRepository(mongodbClient, "Articles")
	getArticles := usecase.NewGetArticlesUseCase(repo)
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		articles, err := getArticles.Execute(c.Context(), repositories.GetArticlesParams{Page: 1, PageSize: 20})
		if err != nil {
			return err
		}
		return c.Render("index", struct {
			Articles  []usecase.ArticleDto
			PageTitle string
		}{
			Articles:  articles,
			PageTitle: "Dev News",
		})
	})

	app.Get("/api/articles", func(c *fiber.Ctx) error {
		c.Context().Logger().Printf("Get articles")
		pageSize := common.ParseInt(c.Query("pageSize"), 10)
		page := common.ParseInt(c.Query("page"), 1)
		log.WithField("pageSize", pageSize).WithField("page", page).Infoln("get articles")
		res, err := getArticles.Execute(c.Context(), repositories.GetArticlesParams{Page: page, PageSize: pageSize})
		if err != nil {
			return err
		}
		return c.Status(200).JSON(res)
	})

	if err := app.Listen(":3000"); err != nil {
		log.WithError(err).Error("can't start server")
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}