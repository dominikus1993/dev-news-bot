using System.Threading.Tasks;
using Discord.Webhook;

namespace DevNews.WebHooks.Application.Services
{
    public class DiscordNetWebHookNotifier : IWebHookNotifier
    {
        private DiscordWebhookClient _client;
        private 

        public DiscordNetWebHookNotifier(DiscordWebhookClient client)
        {
            _client = client;
        }

        public async Task Notify()
        {
            _client.SendMessageAsync()
        }
    }
}