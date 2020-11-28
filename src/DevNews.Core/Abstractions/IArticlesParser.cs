using System.Collections.Generic;
using DevNews.Core.Model;

namespace DevNews.Core.Abstractions
{
    public interface IArticlesParser
    {
        IAsyncEnumerable<Article> Parse();
    }
}