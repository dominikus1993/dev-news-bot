using System;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Repository;
using DevNews.Infrastructure.Persistence.Repository;
using MongoDB.Driver;
using Xunit;

namespace DevNews.WebApp.Tests.Repositories
{
    public class MongoDbFixture : IDisposable
    {
        public IArticlesRepository Repository { get; }

        public MongoDbFixture()
        {
            var client = new MongoClient(Environment.GetEnvironmentVariable("TST_DEVNEWS_DATABASE") ?? "mongodb://root:rootpassword@127.0.0.1:27017");
            Repository = new MongoArticlesRepository(client);
        }


        public void Dispose()
        {
            Console.WriteLine("XDDD");
        }
    }
    
    public class MongoArticlesRepositoryTests : IClassFixture<MongoDbFixture>
    {
        private readonly MongoDbFixture _mongoDbFixture;

        public MongoArticlesRepositoryTests(MongoDbFixture mongoDbFixture)
        {
            _mongoDbFixture = mongoDbFixture;
        }
        
        [Fact]
        public async Task GetArticlesWhenPageHasNegativeValue()
        {
            await Assert.ThrowsAsync<ArgumentOutOfRangeException>(async () =>
                await _mongoDbFixture.Repository.Get(0, 10).ToListAsync());
        }

        [Fact]
        public async Task GetArticlesWhenPageSizeHasNegativeValue()
        {
            await Assert.ThrowsAsync<ArgumentOutOfRangeException>(async () =>
                 await _mongoDbFixture.Repository.Get(1, -10).ToListAsync());
        }
        
        [Fact]
        public async Task GetArticlesWhenPageSizeHasZeroValue()
        {
            await Assert.ThrowsAsync<ArgumentOutOfRangeException>(
                async () =>
                    await _mongoDbFixture.Repository.Get(1, 0).ToListAsync());
        }
        
        
    }
}