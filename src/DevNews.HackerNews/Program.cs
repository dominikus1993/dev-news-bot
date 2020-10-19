using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using HtmlAgilityPack;

namespace DevNews.HackerNews
{
    class Program
    {
        public static async IAsyncEnumerable<string> A()
        {
            var html = new HtmlWeb();
            const string url = "https://news.ycombinator.com/";
            var a = await html.LoadFromWebAsync(url);
            var b = a.DocumentNode.SelectNodes("//*[@class=\"athing\"]");
            b.FindFirst()
            
            yield break;
        }
        static async Task Main(string[] args)
        {
            await A().ToListAsync();
            Console.WriteLine("Hello World!");
        }
    }
}