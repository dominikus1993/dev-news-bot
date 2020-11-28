using DevNews.Core.Repository;
using DevNews.Infrastructure.Persistence.Config;
using DevNews.Infrastructure.Persistence.Repository;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using MongoDB.Driver;

namespace DevNews.Infrastructure.Persistence.DependencyInjection
{
    public static class Extensions
    {
        public static IServiceCollection AddPersistence(this IServiceCollection services, IConfiguration configuration)
        {
            services.AddSingleton<IMongoClient>(_ => new MongoClient(configuration.GetConnectionString("Mongo")));
            services.AddTransient<IArticlesRepository, MongoArticlesRepository>();
            return services;
        }
    }
}