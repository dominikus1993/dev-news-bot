using DevNews.Core.Abstractions;
using DevNews.Infrastructure.Parsers.HackerNews;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Infrastructure.Parsers.DependencyInjection
{
    public static class Extensions
    {
        public static IServiceCollection AddParsers(this IServiceCollection services)
        {
            services.AddTransient<IArticlesParser, HackerNewsArticlesParser>();
            return services;
        }
    }
}