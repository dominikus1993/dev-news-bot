using System;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Reddit;

namespace DevNews.Infrastructure.Parsers.Reddit
{
    internal static class Extensions
    {
        internal static void AddReddit(this IServiceCollection services, IConfiguration configuration)
        {
            services.AddSingleton<RedditConfiguration>(configuration.GetSection("Reddit").Get<RedditConfiguration>());
            services.AddTransient<RedditClient>(sp =>
            {
                var cfg = sp.GetService<RedditConfiguration>() ?? throw new ArgumentNullException(nameof(RedditConfiguration));
                return new RedditClient(cfg.AppId, appSecret: cfg.Secret);
            });
            services.AddTransient<SubRedditParser>();
            services.AddTransient<RedditParser>();
        }
    }
}