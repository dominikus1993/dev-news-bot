using DevNews.Core.Abstractions;
using DevNews.Infrastructure.Notifications.Discord;
using Discord.Webhook;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Infrastructure.Notifications.DependencyInjection
{
    public static class Extensions
    {
        public static IServiceCollection AddNotifiers(this IServiceCollection services, IConfiguration configuration)
        {
            services.AddSingleton(_ => new DiscordWebhookClient(configuration.GetConnectionString("Discord")));
            services.AddTransient<INotifier, DiscordWebHookNotifier>();
            return services;
        }
    }
}