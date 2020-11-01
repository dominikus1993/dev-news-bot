using DevNews.Application.Notifications.Services;
using DevNews.Infrastructure.WebHooks.Services;
using Discord.Webhook;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Infrastructure.WebHooks.DependencyInjection
{
    public static class Extensions
    {
        public static void AddNotifier(this IServiceCollection services, IConfiguration configuration)
        {
            services.AddSingleton(sp => new DiscordWebhookClient(configuration["Discord:WebhookUrl"]));
            services.AddTransient<INotifier, DiscordNetWebHookNotifier>();
            services.AddTransient<DiscordNetWebHookNotifier>();
        }
    }
}