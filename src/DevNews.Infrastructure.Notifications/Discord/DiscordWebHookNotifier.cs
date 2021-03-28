using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using Discord;
using Discord.Webhook;
using Microsoft.Extensions.Logging;

namespace DevNews.Infrastructure.Notifications.Discord
{
    internal class DiscordWebHookNotifier : INotifier
    {
        private readonly DiscordWebhookClient _discordWebhookClient;
        private readonly ILogger<DiscordWebHookNotifier> _logger;

        public DiscordWebHookNotifier(DiscordWebhookClient discordWebhookClient, ILogger<DiscordWebHookNotifier> logger)
        {
            _discordWebhookClient = discordWebhookClient;
            _logger = logger;
        }

        public async Task Notify(IEnumerable<Article> articles)
        {
            _logger.LogInformation("Start Sending articles");   
            var embeds = articles
                .Select(article => article.CreateEmbed())
                .ToList();
            await _discordWebhookClient.SendMessageAsync("Witam serdecznie, oto nowe newsy", embeds: embeds);
            _logger.LogInformation("Send articles succeed"); 
        }
    }
}