using System.Collections.Generic;
using DevNews.HackerNews.Application.Dto;
using DevNews.HackerNews.Application.Servies;

namespace DevNews.HackerNews.Application.UseCases
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