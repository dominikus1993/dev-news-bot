using DevNews.Application.HackerNews.Servies;
using DevNews.Application.HackerNews.UseCases;
using DevNews.Core.HackerNews;
using DevNews.Infrastructure.HackerNews.Repositories;
using DevNews.Infrastructure.HackerNews.Services;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Infrastructure.HackerNews.DependencyInjection
{
    public static class Extensions
    {
        public static void AddHackerNews(this IServiceCollection services)
        {
            services.AddTransient<IHackerNewsRepository, FakeHackerNewsRepository>();
            services.AddTransient<IHackerNewsParser, HtmlHackerNewsParser>();
        }
    }
}