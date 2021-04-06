using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Model;
using DevNews.Core.Repository;
using DevNews.Infrastructure.Persistence.Model;
using LanguageExt;
using MongoDB.Driver;
using MongoDB.Driver.Linq;
using static LanguageExt.Prelude;

namespace DevNews.Infrastructure.Persistence.Repository
{
    public class MongoArticlesRepository : IArticlesRepository
    {
        private readonly IMongoCollection<MongoArticle> _articles;

        public MongoArticlesRepository(IMongoClient client)
        {
            _articles = GetCollection(GetDatabase(client));
        }

        public async Task<bool> Exists(Article article) =>
            await _articles.AsQueryable().AnyAsync(x => x.Title == article.Title);

        public async Task<Either<Exception, Unit>> InsertMany(IEnumerable<Article> articles)
        {
            try
            {
                var writes = articles
                    .Select(article => new MongoArticle(article))
                    .Select(article => new InsertOneModel<MongoArticle>(article))
                    .ToList();

                if (writes.Any())
                {
                    await _articles.BulkWriteAsync(writes);
                }

                return Right(Unit.Default);
            }
            catch (Exception e)
            {
                return Left(e);
            }
        }

        private static IMongoDatabase GetDatabase(IMongoClient client)
        {
            return client.GetDatabase("Articles");
        }

        private static IMongoCollection<MongoArticle> GetCollection(IMongoDatabase db)
        {
            return db.GetCollection<MongoArticle>("articles");
        }
    }
}