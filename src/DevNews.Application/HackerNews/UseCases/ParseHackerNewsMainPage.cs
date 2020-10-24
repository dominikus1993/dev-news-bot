using System.Collections.Generic;
using DevNews.Application.HackerNews.Dto;
using DevNews.Application.HackerNews.Servies;

namespace DevNews.Application.HackerNews.UseCases
{
    public class ParseHackerNewsMainPageUseCase
    {
        private IHackerNewsParser _hackerNewsParser;

        public ParseHackerNewsMainPageUseCase(IHackerNewsParser hackerNewsParser)
        {
            _hackerNewsParser = hackerNewsParser;
        }

        public IAsyncEnumerable<ArticleDto> Execute()
        {
            return _hackerNewsParser.Parse();
        }
    }
}