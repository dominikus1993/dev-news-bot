using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using Discord;
using Discord.Webhook;

namespace DevNews.Infrastructure.Notifications.Discord
{
    internal class DiscordWebHookNotifier : INotifier
    {
        private readonly DiscordWebhookClient _discordWebhookClient;

        public DiscordWebHookNotifier(DiscordWebhookClient discordWebhookClient)
        {
            _discordWebhookClient = discordWebhookClient;
        }

        public async Task Notify(IEnumerable<Article> articles)
        {
            var embeds = articles
                .Select(article => article.CreateEmbed())
                .ToList();
            await _discordWebhookClient.SendMessageAsync("Witam serdecznie, oto nowe newsy", embeds: embeds);
        }
    }
}