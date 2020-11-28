using DevNews.Core.Repository;
using DevNews.Infrastructure.Persistence.Config;
using DevNews.Infrastructure.Persistence.Repository;
using Microsoft.Extensions.DependencyInjection;
using MongoDB.Driver;

namespace DevNews.Infrastructure.Persistence.DependencyInjection
{
    public static class Extensions
    {
        public static IServiceCollection AddPersistence(this IServiceCollection services, PersistenceConfiguration configuration)
        {
            services.AddSingleton<IMongoClient>(sp => new MongoClient(configuration.MongoConnectionString));
            services.AddTransient<IArticlesRepository, MongoArticlesRepository>();
            return services;
        }
    }
}