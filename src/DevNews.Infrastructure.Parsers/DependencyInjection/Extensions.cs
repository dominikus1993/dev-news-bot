using System;
using DevNews.Core.Abstractions;
using DevNews.Infrastructure.Parsers.Dotnetomaniak;
using DevNews.Infrastructure.Parsers.HackerNews;
using DevNews.Infrastructure.Parsers.Reddit;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Reddit;

namespace DevNews.Infrastructure.Parsers.DependencyInjection
{
    public static class Extensions
    {
        public static IServiceCollection AddParsers(this IServiceCollection services, IConfiguration configuration)
        {
            services.AddTransient<IArticlesParser, HackerNewsArticlesParser>();
            services.AddTransient<IArticlesParser, DotnetomaniakArticlesParser>();
            services.AddReddit(configuration);
            return services;
        }
        
    }
}