using System.Linq;
using System.Threading.Tasks;
using DevNews.Infrastructure.Parsers.Dotnetomaniak;
using FluentAssertions;
using Xunit;

namespace DevNews.Infrastructure.Parsers.UnitTests.Parsers
{
    public class DotnetomaniakArticleParserTests
    {
        [Fact]
        public async Task TestParsing()
        {
            var parser = new DotnetomaniakArticlesParser();
            var subject = await parser.Parse().ToListAsync();
            subject.Should().NotBeEmpty();
        }
    }
}