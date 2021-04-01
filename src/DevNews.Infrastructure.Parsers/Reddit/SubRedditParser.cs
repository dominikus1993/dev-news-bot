using System;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Model;
using HtmlAgilityPack;

namespace DevNews.Infrastructure.Parsers.Reddit
{
    public class SubRedditParser
    {

        public Task<Article[]> Parse(string name)
        {
            var url = $"https://www.reddit.com/r/{name}";
            var html = new HtmlWeb();
            var document = await html.LoadFromWebAsync(url);
            
            var nodes = document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]")
                .Select(static node => new Article(node.InnerText, node.GetAttributeValue("href", null)));
            
        }
    }
}