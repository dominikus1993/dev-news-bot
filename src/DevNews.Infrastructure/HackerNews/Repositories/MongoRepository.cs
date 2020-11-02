using System.Collections.Generic;
using DevNews.Core.HackerNews;
using DevNews.Core.Model;

namespace DevNews.Infrastructure.HackerNews.Repositories
{
    public class MongoHackerNewsRepository : IHackerNewsRepository
    {
        public IAsyncEnumerable<(Article, bool)> Exists(IEnumerable<Article> articles)
        {
            throw new System.NotImplementedException();
        }
    }
}