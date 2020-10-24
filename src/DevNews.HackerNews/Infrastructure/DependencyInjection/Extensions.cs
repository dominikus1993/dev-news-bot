using DevNews.HackerNews.Application.Servies;
using DevNews.HackerNews.Application.UseCases;
using DevNews.HackerNews.Infrastructure.Services;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.HackerNews.Infrastructure.DependencyInjection
{
    public static class Extensions
    {
        public static void AddHackerNews(this IServiceCollection services)
        {
            services.AddTransient<IHackerNewsParser, HtmlHackerNewsParser>();
            services.AddScoped<ParseHackerNewsMainPageUseCase>();
        }
    }
}