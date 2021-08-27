using System.Collections.Generic;
using System.Threading;
using DevNews.Core.Model;

namespace DevNews.Core.Abstractions
{
    public interface IArticlesProvider
    {
        IAsyncEnumerable<Article> Provide(CancellationToken cancellationToken = default);
    }
}