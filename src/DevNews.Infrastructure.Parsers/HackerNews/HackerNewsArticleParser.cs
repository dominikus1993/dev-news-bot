using System.Collections.Generic;
using System.Runtime.CompilerServices;
using System.Threading;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using HtmlAgilityPack;

[assembly: InternalsVisibleTo("DevNews.Infrastructure.Parsers.UnitTests")]
namespace DevNews.Infrastructure.Parsers.HackerNews
{
    internal class HackerNewsArticlesParser : IArticlesParser
    {
        private const string HackerNewsUrl = "https://news.ycombinator.com/";
        
        public async IAsyncEnumerable<Article> Parse([EnumeratorCancellation] CancellationToken cancellationToken = default)
        {
            var html = new HtmlWeb();
            var document = await html.LoadFromWebAsync(HackerNewsUrl, cancellationToken);

            var nodes = document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]");

            foreach (var node in nodes)
            {
                yield return new Article(node.InnerText, node.GetAttributeValue("href", null));
            }
        }
    }
}