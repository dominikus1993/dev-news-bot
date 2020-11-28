using System.Threading.Tasks;
using Cocona;
using DevNews.Cli.Infrastructure;
using DevNews.Core.DependencyInjection;
using DevNews.Core.UseCases;
using DevNews.Infrastructure.Notifications.DependencyInjection;
using DevNews.Infrastructure.Parsers.DependencyInjection;
using DevNews.Infrastructure.Persistence.DependencyInjection;
using Microsoft.Extensions.Configuration;

namespace DevNews.Cli
{
    class Program
    {
        static async Task Main(string[] args) =>
            await CoconaApp.Create()
                .UseLogger("DevNews.Cli")
                .ConfigureServices((ctx, services) =>
                {
                    services.AddPersistence(ctx.Configuration);
                    services.AddNotifiers(ctx.Configuration);
                    services.AddParsers();
                    services.AddCore();
                })
                .ConfigureAppConfiguration(builder =>
                {
                    builder
                        .AddJsonFile("appsettings.json", true, true)
                        .AddEnvironmentVariables();
                })
                .RunAsync<Program>(args);

        public async Task
            ProduceNews([FromService] ParseArticlesAndSendItUseCase parseArticlesAndSendItUseCase) =>
            await parseArticlesAndSendItUseCase.Execute(new ParseArticlesAndSendItParam(5));
    }
}