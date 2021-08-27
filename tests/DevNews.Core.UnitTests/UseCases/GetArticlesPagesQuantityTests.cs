using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using DevNews.Core.Model;
using DevNews.Core.Repository;
using DevNews.Core.UseCases;
using FluentAssertions;
using Moq;
using Xunit;

namespace DevNews.Core.UnitTests.UseCases
{
    public class GetArticlesPagesQuantityTests
    {
        [Fact]
        public async Task GetArticlesWhenPageHasZeroValue()
        {
            // Arrange 
            var repo = new Mock<IArticlesRepository>();
            var useCase = new GetArticlesPagesQuantity(repo.Object);
            
            // Test
            await Assert.ThrowsAsync<ArgumentOutOfRangeException>(() =>
                useCase.Execute(0));
        }
        
        [Fact]
        public async Task GetArticlesWhenPageHasNegativeValue()
        {
            // Arrange 
            var repo = new Mock<IArticlesRepository>();
            var useCase = new GetArticlesPagesQuantity(repo.Object);
            
            // Test
            await Assert.ThrowsAsync<ArgumentOutOfRangeException>(() =>
                useCase.Execute(-1));
        }
        
        [Theory]
        [InlineData(10, 100, 10)]
        [InlineData(10, 95, 10)]
        [InlineData(5, 5, 1)]
        [InlineData(5, 0, 0)]
        [InlineData(5, 4, 1)]
        [InlineData(10, 10, 1)]
        public async Task GetArticlesPages(int pageSize, int articlesQuantity, int expectedPagesQuantity)
        {
            // Arrange 
            var repo = new Mock<IArticlesRepository>();
            repo.Setup(repository => repository.Count(It.IsAny<CancellationToken>())).ReturnsAsync(articlesQuantity);
            var useCase = new GetArticlesPagesQuantity(repo.Object);
            
            // Test
            var subject = await useCase.Execute(pageSize);

            subject.Should().Be(expectedPagesQuantity);
        }
    }
}