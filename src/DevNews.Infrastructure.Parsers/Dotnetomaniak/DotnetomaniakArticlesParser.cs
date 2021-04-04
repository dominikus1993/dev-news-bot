using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using HtmlAgilityPack;

[assembly: InternalsVisibleTo("DevNews.Infrastructure.Parsers.UnitTests")]
namespace DevNews.Infrastructure.Parsers.Dotnetomaniak
{
    internal class DotnetomaniakArticlesParser : IArticlesParser
    {
        private const string DotnetoManiakUrl = "https://dotnetomaniak.pl/";
        private static readonly Uri DotnetoManiakUri = new(DotnetoManiakUrl);
        
        public async IAsyncEnumerable<Article> Parse()
        {
            var html = new HtmlWeb();
            var document = await html.LoadFromWebAsync(DotnetoManiakUrl);

            var nodes = document.DocumentNode.SelectNodes("//*[@class=\"article\"]")
                .Select(x => x.ChildNodes);

            foreach (var node in nodes)
            {
                yield return CreateArticle(node);
            }
        }

        private static Article CreateArticle(HtmlNodeCollection nodes)
        {
            var titleNode = nodes.First(div => div.HasClass("title"));
            var href = titleNode.ChildNodes.FindFirst("a").GetAttributeValue("href", null);
            var link = new Uri(DotnetoManiakUri, href).AbsoluteUri;
            var description = nodes.First(div => div.HasClass("description")).ChildNodes.FindFirst("span").InnerText;
            var title = titleNode.InnerText;
            return new Article(title, description, link);
        }
    }
}