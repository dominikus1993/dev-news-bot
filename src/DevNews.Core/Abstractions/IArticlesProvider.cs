using System.Collections.Generic;
using DevNews.Core.Model;

namespace DevNews.Core.Abstractions
{
    public interface IArticlesProvider
    {
        IAsyncEnumerable<Article> Provide();
    }
}