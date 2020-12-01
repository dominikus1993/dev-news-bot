using System.Collections.Generic;
using System.Linq;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using HtmlAgilityPack;

namespace DevNews.Infrastructure.Parsers.Dotnetomaniak
{
    public class DotnetomaniakArticlesParser : IArticlesParser
    {
        public const string DotnetoManiakUrl = "https://dotnetomaniak.pl/";
        public async IAsyncEnumerable<Article> Parse()
        {
            var html = new HtmlWeb();
            var document = await html.LoadFromWebAsync(DotnetoManiakUrl);
            
            var nodes = document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]")
                .Select(node => new Article(node.InnerText, node.GetAttributeValue("href", null)));
            
            foreach (var article in nodes)
            {
                yield return article;
            }
        }
    }
}