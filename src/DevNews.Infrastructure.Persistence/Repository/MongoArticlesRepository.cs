using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Threading;
using System.Threading.Tasks;
using DevNews.Core.Model;
using DevNews.Core.Repository;
using DevNews.Infrastructure.Persistence.Model;
using LanguageExt;
using MongoDB.Driver;
using MongoDB.Driver.Linq;
using static LanguageExt.Prelude;

[assembly: InternalsVisibleTo("DevNews.Persistence.UnitTests")]
namespace DevNews.Infrastructure.Persistence.Repository
{
    internal class MongoArticlesRepository : IArticlesRepository
    {
        private readonly IMongoCollection<MongoArticle> _articles;
        public const string ArticlesDatabase = "Articles";
        public const string ArticlesCollection = "articles";

        public MongoArticlesRepository(IMongoClient client)
        {
            _articles = GetCollection(GetDatabase(client));
        }

        public async Task<bool> Exists(Article article, CancellationToken cancellationToken = default) =>
            await _articles.AsQueryable().AnyAsync(x => x.Title == article.Title, cancellationToken: cancellationToken);

        public async Task<Either<Exception, Unit>> InsertMany(IEnumerable<Article> articles, CancellationToken cancellationToken = default)
        {
            try
            {
                var writes = articles
                    .Select(article => new MongoArticle(article))
                    .Select(article => new InsertOneModel<MongoArticle>(article))
                    .ToList();

                if (writes.Count > 0)
                {
                    await _articles.BulkWriteAsync(writes, cancellationToken: cancellationToken);
                }

                return Right(Unit.Default);
            }
            catch (Exception e)
            {
                return Left(e);
            }
        }

        public async IAsyncEnumerable<Article> Get(int page, int pageSize, [EnumeratorCancellation]CancellationToken cancellationToken = default)
        {
            if (page <= 0)
            {
                throw new ArgumentOutOfRangeException(nameof(page));
            }

            if (pageSize <= 0)
            {
                throw new ArgumentOutOfRangeException(nameof(pageSize));
            }

            var skip = (page - 1) * pageSize;

            var result = await _articles
                .AsQueryable()
                .OrderBy(x => x.CrawledAt)
                .Skip(skip)
                .Take(pageSize)
                .ToListAsync(cancellationToken: cancellationToken);

            if (result is null)
            {
                yield break;
            }

            foreach (var mongoArticle in result)
            {
                yield return mongoArticle.AsArticle();
            }
        }

        public async Task<long> Count(CancellationToken cancellationToken = default) => 
            await _articles.CountDocumentsAsync(FilterDefinition<MongoArticle>.Empty, cancellationToken: cancellationToken);

        public static IMongoDatabase GetDatabase(IMongoClient client)
        {
            return client.GetDatabase(ArticlesDatabase);
        }

        public static IMongoCollection<MongoArticle> GetCollection(IMongoDatabase db)
        {
            return db.GetCollection<MongoArticle>(ArticlesCollection);
        }
    }
}