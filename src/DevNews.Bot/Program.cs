using System.Threading.Tasks;
using Akka.Actor;
using DevNews.Application.Framework.DependencyInjection;
using DevNews.Infrastructure.HackerNews.DependencyInjection;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace DevNews.DiscordBot
{
    public class Program
    {
        public static async Task Main(string[] args)
        {
            await CreateHostBuilder(args).Build().RunAsync();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureServices((hostContext, services) =>
                {
                    services.AddApplication();
                    services.AddHackerNews();
                    services.AddHostedService<Worker>();
                });
    }
}
