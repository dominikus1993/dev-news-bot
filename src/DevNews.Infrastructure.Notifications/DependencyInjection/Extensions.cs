using System;
using DevNews.Core.Abstractions;
using DevNews.Infrastructure.Notifications.Discord;
using DevNews.Infrastructure.Notifications.MicrosoftTeams;
using DevNews.Infrastructure.Notifications.MicrosoftTeams.Model;
using Discord.Webhook;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Infrastructure.Notifications.DependencyInjection
{
    public static class Extensions
    {
        public static IServiceCollection AddNotifiers(this IServiceCollection services, IConfiguration configuration)
        {
            AddDiscord(services, configuration);
            AddMicrosoftTeams(services, configuration);
            return services;
        }

        private static IServiceCollection AddDiscord(IServiceCollection services, IConfiguration configuration)
        {
            services.AddSingleton(_ => new DiscordWebhookClient(configuration.GetConnectionString("Discord")));
            services.AddTransient<INotifier, DiscordWebHookNotifier>();
            return services;
        }
        
        private static IServiceCollection AddMicrosoftTeams(IServiceCollection services, IConfiguration configuration)
        {
            if (bool.TryParse(configuration["MicrosoftTeams:Enabled"], out var enabled) && enabled)
            {
                services.AddScoped<ITeamsMessageSerializer, JsonTeamsMessageSerializer>();
                services.AddTransient<INotifier, MicrosoftTeamsWebHookNotifier>();
                services.AddHttpClient(MicrosoftTeamsWebHookNotifier.MicrosoftTeamsWebHookNotifierApi, client =>
                {
                    client.BaseAddress = new Uri(configuration.GetConnectionString("Teams") ?? throw new ArgumentNullException());
                });
            }

            return services;
        }
    }
}