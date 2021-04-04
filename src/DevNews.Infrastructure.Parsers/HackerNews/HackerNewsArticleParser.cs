using System.Collections.Generic;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using HtmlAgilityPack;

namespace DevNews.Infrastructure.Parsers.HackerNews
{
    public class HackerNewsArticlesParser : IArticlesParser
    {
        private const string HackerNewsUrl = "https://news.ycombinator.com/";
        
        public async IAsyncEnumerable<Article> Parse()
        {
            var html = new HtmlWeb();
            var document = await html.LoadFromWebAsync(HackerNewsUrl);

            var nodes = document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]");

            foreach (var node in nodes)
            {
                yield return new Article(node.InnerText, node.GetAttributeValue("href", null));
            }
        }
    }
}