using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Parsers.HackerNews;
using FluentAssertions;
using Xunit;

namespace DevNews.Core.UnitTests.Parsers
{
    public class HackerNewsParseTests
    {
        [Fact]
        public async Task TestParsing()
        {
            var parser = new HackerNewsArticlesParser();
            var subject = await parser.Parse().ToListAsync();
            subject.Should().NotBeEmpty();
        }
    }
}