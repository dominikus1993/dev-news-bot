using System.Collections.Generic;
using DevNews.HackerNews.Application.Dto;

namespace DevNews.HackerNews.Application.Servies
{
    public interface IHackerNewsParser
    {
        IAsyncEnumerable<ArticleDto> Parse();
    }
}