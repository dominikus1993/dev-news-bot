using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Xml;
using DevNews.Shared.Messages;
using Discord;
using Discord.Webhook;

namespace DevNews.WebHooks.Application.Services
{
    public class DiscordNetWebHookNotifier : IWebHookNotifier
    {
        private DiscordWebhookClient _client;

        public DiscordNetWebHookNotifier(DiscordWebhookClient client)
        {
            _client = client;
        }

        public async Task Notify(IEnumerable<Article> articles)
        {
            var embeds = articles.Select(article => new EmbedBuilder().WithUrl(article.Url).WithTitle(article.Title))
                .Select(x => x.Build()).ToList();
            await _client.SendMessageAsync("Nowe newsy od HackerNews", false, embeds);
        }
    }
}