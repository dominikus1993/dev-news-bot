using System;
using System.Linq;
using System.Net.Http;
using System.Threading.Tasks;
using DevNews.Infrastructure.Parsers.HackerNews;
using DevNews.Infrastructure.Parsers.Reddit;
using FluentAssertions;
using Xunit;

namespace DevNews.Infrastructure.Parsers.UnitTests.Parsers
{
    public class RedditParserTests
    {
        [Fact]
        public async Task TestParsing()
        {
            using var httpclient = new HttpClient()
            {
                BaseAddress = new Uri("https://www.reddit.com/")
            };
            var parser = new RedditParser(new RedditConfiguration() { SubReddits = new []{ "dotnet" }}, new SubRedditParser(httpclient));
            var subject = await parser.Parse().ToListAsync();
            subject.Should().NotBeEmpty();
        }
    }
}