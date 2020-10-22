using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using DevNews.WebHooks.Application.Services;
using Discord.Webhook;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace DevNews.DiscordBot
{
    public class Program
    {
        public static async Task Main(string[] args)
        {
            using var client = new DiscordWebhookClient(string.Empty);
            var service = new DiscordNetWebHookNotifier(client);
            await service.Notify();
            //CreateHostBuilder(args).Build().Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureServices((hostContext, services) =>
                {
                    services.AddHostedService<Worker>();
                });
    }
}
