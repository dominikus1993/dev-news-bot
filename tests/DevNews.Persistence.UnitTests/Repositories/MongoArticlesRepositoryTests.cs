using System;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Model;
using DevNews.Core.Repository;
using DevNews.Infrastructure.Persistence.Repository;
using FluentAssertions;
using MongoDB.Driver;
using Xunit;

namespace DevNews.WebApp.Tests.Repositories
{
    public class MongoDbFixture : IDisposable
    {
        public IArticlesRepository Repository { get; }
        public IMongoClient Client { get; }

        public MongoDbFixture()
        {
            Client = new MongoClient(Environment.GetEnvironmentVariable("TST_DEVNEWS_DATABASE") ?? "mongodb://root:rootpassword@127.0.0.1:27017");
            Repository = new MongoArticlesRepository(Client);
        }


        public void Dispose()
        {
            Client.DropDatabase(MongoArticlesRepository.ArticlesDatabase);
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
        
        
        [Fact]
        public async Task GetArticlesWhenCollectionIsEmpty()
        {
            await MongoArticlesRepository.GetDatabase(_mongoDbFixture.Client)
                .DropCollectionAsync(MongoArticlesRepository.ArticlesCollection);

            var result = await _mongoDbFixture.Repository.Get(1, 10).ToListAsync();

            result.Should().BeEmpty();
        }
        
        [Fact]
        public async Task GetArticlesWhenCollectionIsNotEmpty()
        {
            
            await MongoArticlesRepository.GetDatabase(_mongoDbFixture.Client)
                .DropCollectionAsync(MongoArticlesRepository.ArticlesCollection);

            await _mongoDbFixture.Repository.InsertMany(Enumerable.Range(1, 20)
                .Select(x => new Article($"xDD {x}", "xD", $"http://www.xD.com/{x}")));

            // TEST
            var subject = await _mongoDbFixture.Repository.Get(1, 10).ToListAsync();

            subject.Should().HaveCount(10);
            
            subject = await _mongoDbFixture.Repository.Get(2, 10).ToListAsync();
            
            subject.Should().HaveCount(10);
        }
        
        
        [Fact]
        public async Task CheckExistenceWhenArticleNotExists()
        {
            
            await MongoArticlesRepository.GetDatabase(_mongoDbFixture.Client)
                .DropCollectionAsync(MongoArticlesRepository.ArticlesCollection);

            await _mongoDbFixture.Repository.InsertMany(Enumerable.Range(1, 20)
                .Select(x => new Article($"xDD {x}", "xD", $"http://www.xD.com/{x}")));

            // TEST
            var subject = await _mongoDbFixture.Repository.Exists(new Article("xDD", "www.xD.com"));

            subject.Should().BeFalse();
        }
        
        
        [Fact]
        public async Task CheckExistenceWhenArticleExists()
        {
            
            await MongoArticlesRepository.GetDatabase(_mongoDbFixture.Client)
                .DropCollectionAsync(MongoArticlesRepository.ArticlesCollection);

            await _mongoDbFixture.Repository.InsertMany(Enumerable.Range(1, 20)
                .Select(x => new Article($"xDD {x}", "xD", $"http://www.xD.com/{x}")));

            // TEST
            var subject = await _mongoDbFixture.Repository.Exists(new Article("xDD 1", "xD","http://www.xD.com/1"));

            subject.Should().BeTrue();
        }
    }
}