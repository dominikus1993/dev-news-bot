using System;
using System.Threading.Tasks;
using DevNews.Core.Model;

namespace DevNews.Infrastructure.Parsers.Reddit
{
    public class SubRedditParser
    {

        public Task<Article[]> Parse(string name)
        {
            return Task.FromResult(Array.Empty<Article>());
        }
    }
}