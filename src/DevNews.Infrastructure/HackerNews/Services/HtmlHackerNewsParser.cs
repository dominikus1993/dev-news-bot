using System.Collections.Generic;
using System.Linq;
using DevNews.Application.HackerNews.Dto;
using DevNews.Application.HackerNews.Servies;
using HtmlAgilityPack;

namespace DevNews.Infrastructure.HackerNews.Services
{
    public class HtmlHackerNewsParser : IHackerNewsParser
    {
        private const string HackerNewsUrl = "https://news.ycombinator.com/";
        public async IAsyncEnumerable<ArticleDto> Parse()
        {
            var html = new HtmlWeb();
            var document = await html.LoadFromWebAsync(HackerNewsUrl);
            var nodes =
                document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]")
                    .Select(e => (link: e.GetAttributeValue("href", null), title: e.InnerText));

            foreach (var (link, title) in nodes)
            {
                yield return new ArticleDto(title, link);
            }
        }
    }
}