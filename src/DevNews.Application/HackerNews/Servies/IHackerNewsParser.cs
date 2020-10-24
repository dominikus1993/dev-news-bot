using System.Collections.Generic;
using DevNews.Application.HackerNews.Dto;

namespace DevNews.Application.HackerNews.Servies
{
    public interface IHackerNewsParser
    {
        IAsyncEnumerable<ArticleDto> Parse();
    }
}