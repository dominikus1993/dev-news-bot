using System;
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
        private static readonly Uri DotnetoManiakUri = new(DotnetoManiakUrl);
        public async IAsyncEnumerable<Article> Parse()
        {
            var html = new HtmlWeb();
            var document = await html.LoadFromWebAsync(DotnetoManiakUrl);

            var nodes = document.DocumentNode.SelectNodes("//*[@class=\"article\"]")
                .Select(x => x.ChildNodes).Select(static n => CreateArticle(n));

            foreach (var article in nodes)
            {
                yield return article;
            }
        }

        private static Article CreateArticle(HtmlNodeCollection nodes)
        {
            var titleNode = nodes.First(div => div.HasClass("title"));
            var href = titleNode.ChildNodes.FindFirst("a").GetAttributeValue("href", null);
            var link = new Uri(DotnetoManiakUri, href);
            var descElement = nodes.First(div => div.HasClass("description"));
            var description = descElement.ChildNodes.FindFirst("span").InnerText;
            return new Article(titleNode.InnerHtml, description, link.AbsoluteUri);
        }
    }
}