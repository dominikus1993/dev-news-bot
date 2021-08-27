using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using DevNews.Core.Model;
using DevNews.Core.Repository;
using DevNews.Core.UseCases;
using Moq;
using Xunit;

namespace DevNews.Core.UnitTests.UseCases
{
    public class GetArticlesUseCaseFixture : IDisposable
    {
        public GetArticlesUseCaseFixture()
        {
            var article = new Article("test", "https://scienceintegritydigest.com/2020/12/20/paper-about-herbalife-related-patient-death-removed-after-company-threatens-to-sue-the-journal/");
            var mock = new Mock<IArticlesRepository>();
            mock.Setup(x => x.Get(It.IsAny<int>(), It.IsAny<int>(), It.IsAny<CancellationToken>()))
                .Returns<int, int>((page, pageSize) => new[] {article}.ToAsyncEnumerable());
            MockRepository = mock;
            GetArticles = new GetArticles(mock.Object);
        }

        public void Dispose()
        {
            // ... clean up test data from the database ...
        }

        public GetArticles GetArticles { get; set; }
        public Mock<IArticlesRepository> MockRepository { get; set; }
    }

    public class GetArticlesTests : IClassFixture<GetArticlesUseCaseFixture>
    {
        GetArticlesUseCaseFixture _fixture;

        public GetArticlesTests(GetArticlesUseCaseFixture fixture)
        {
            this._fixture = fixture;
        }

        [Fact]
        public async Task GetArticlesWhenPageHasNegativeValue()
        {
            await Assert.ThrowsAsync<ArgumentOutOfRangeException>(() =>
                _fixture.GetArticles.Execute(new GetArticlesQuery(-1, 10)));
        }

        [Fact]
        public async Task GetArticlesWhenPageSizeHasNegativeValue()
        {
            await Assert.ThrowsAsync<ArgumentOutOfRangeException>(() =>
                _fixture.GetArticles.Execute(new GetArticlesQuery(1, -10)));
        }
        
        [Fact]
        public async Task GetArticlesWhenPageSizeHasZeroValue()
        {
            await Assert.ThrowsAsync<ArgumentOutOfRangeException>(() =>
                _fixture.GetArticles.Execute(new GetArticlesQuery(1, 0)));
        }

        [Fact]
        public async Task GetArticles()
        {
            var result = await _fixture.GetArticles.Execute(new GetArticlesQuery(1, 10));
            _fixture.MockRepository.Verify(repository => repository.Get(1, 10, It.IsAny<CancellationToken>()), Times.Once);
        }
    }
}