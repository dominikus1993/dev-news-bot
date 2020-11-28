using DevNews.Core.Abstractions;
using DevNews.Core.Notifications;
using DevNews.Core.Providers;
using DevNews.Core.UseCases;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Core.DependencyInjection
{
    public static class Extensions
    {
        public static IServiceCollection AddParsers(this IServiceCollection services)
        {
            services.AddTransient<IArticlesProvider, ChannelArticlesProvider>();
            services.AddTransient<ParseArticlesAndSendItUseCase>();
            services.AddTransient<INotificationBroadcaster, ChannelsNotificationBroadcaster>();
            return services;
        }
    }
}