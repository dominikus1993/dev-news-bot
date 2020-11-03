using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Core.HackerNews;
using DevNews.Core.Model;
using MongoDB.Driver;

namespace DevNews.Infrastructure.HackerNews.Repositories
{
    public class MongoHackerNewsRepository : IHackerNewsRepository
    {
        private MongoClient _client;

        public MongoHackerNewsRepository(MongoClient client)
        {
            _client = client;
        }

        public IAsyncEnumerable<(Article, bool)> Exists(IEnumerable<Article> articles)
        {
            throw new System.NotImplementedException();
        }

        public async Task InsertMany(IEnumerable<Article> articles)
        {
            throw new System.NotImplementedException();
        }
    }
}