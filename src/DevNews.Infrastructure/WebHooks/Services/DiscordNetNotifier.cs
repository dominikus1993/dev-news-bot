using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Application.Notifications.Services;
using Discord;
using Discord.Webhook;
using DevNews.Core.HackerNews;
using  DevNews.Core.Model;
namespace DevNews.Infrastructure.WebHooks.Services
{
    public class DiscordNetWebHookNotifier : INotifier
    {
        private DiscordWebhookClient _client;

        public DiscordNetWebHookNotifier(DiscordWebhookClient client)
        {
            _client = client;
        }

        public async Task Notify(IEnumerable<Article> articles)
        {
            var embeds = articles.Select(article => new EmbedBuilder().WithUrl(article.Link).WithTitle(article.Title))
                .Select(x => x.Build()).ToList();
            await _client.SendMessageAsync("Nowe newsy od HackerNews", false, embeds);
        }
    }
}