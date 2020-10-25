using DevNews.Application.HackerNews.Servies;
using DevNews.Application.HackerNews.UseCases;
using DevNews.Infrastructure.HackerNews.Services;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Infrastructure.HackerNews.DependencyInjection
{
    public static class Extensions
    {
        public static void AddHackerNews(this IServiceCollection services)
        {
            services.AddTransient<IHackerNewsParser, HtmlHackerNewsParser>();
        }
    }
}