using System;
using DevNews.Core.Abstractions;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Infrastructure.Parsers.Reddit
{
    internal static class Extensions
    {
        internal static void AddReddit(this IServiceCollection services, IConfiguration configuration)
        {
            services.AddSingleton(configuration.GetSection("Reddit").Get<RedditConfiguration>());
            services.AddHttpClient<SubRedditParser>(client =>
            {
                client.BaseAddress = new Uri("https://www.reddit.com/");
            });
            services.AddTransient<IArticlesParser, RedditParser>();
        }
    }
}