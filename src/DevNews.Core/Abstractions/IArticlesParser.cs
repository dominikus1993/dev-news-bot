using System.Collections.Generic;
using System.Runtime.CompilerServices;
using System.Threading;
using DevNews.Core.Model;

namespace DevNews.Core.Abstractions
{
    public interface IArticlesParser
    {
        IAsyncEnumerable<Article> Parse(CancellationToken cancellationToken = default);
    }
}