using System.Collections.Generic;
using System.Linq;
using DevNews.Core.HackerNews;
using DevNews.Core.Model;

namespace DevNews.Infrastructure.HackerNews.Repositories
{
    public class FakeHackerNewsRepository : IHackerNewsRepository
    {
        public IAsyncEnumerable<(Article, bool)> Exists(IEnumerable<Article> articles)
        {
            return articles.Select(x => (x, false)).ToAsyncEnumerable();
        }
    }
}