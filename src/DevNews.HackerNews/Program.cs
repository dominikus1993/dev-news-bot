using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using HtmlAgilityPack;

namespace DevNews.HackerNews
{
    class Article
    {
        public Article(string title, string link)
        {
            Title = title;
            Link = link;
        }

        public string Title { get; }
        public string Link { get; }

        public override string ToString()
        {
            return $"{nameof(Title)}: {Title}, {nameof(Link)}: {Link}";
        }
    }
    class Program
    {
        public static async IAsyncEnumerable<Article> A()
        {
            var html = new HtmlWeb();
            const string url = "https://news.ycombinator.com/";
            var document = await html.LoadFromWebAsync(url);
            var nodes =
                document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]")
                    .Select(e => (link: e.GetAttributeValue("href", null), title: e.InnerText));

            foreach (var (link, title) in nodes)
            {
                yield return new Article(title, link);
            }
            
        }
        static async Task Main(string[] args)
        {
            await foreach (var article in A())
            {
                Console.WriteLine(article);
            }
            
            Console.WriteLine("Hello World!");
        }
    }
}