using System.Threading.Tasks;
using Cocona;
using DevNews.Cli.Infrastructure;
using DevNews.Core.DependencyInjection;
using DevNews.Core.UseCases;
using DevNews.Infrastructure.Notifications.DependencyInjection;
using DevNews.Infrastructure.Parsers.DependencyInjection;
using DevNews.Infrastructure.Persistence.DependencyInjection;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;

namespace DevNews.Cli
{
    public class Program
    {
        private ILogger<Program> _logger;

        public Program(ILogger<Program> logger)
        {
            _logger = logger;
        }

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
            ProduceNews(int articleQuantity, [FromService] ParseArticlesAndSendItUseCase parseArticlesAndSendItUseCase)
        {
            _logger.LogInformation("Start Producing News");
            await parseArticlesAndSendItUseCase.Execute(new ParseArticlesAndSendItParam(articleQuantity));
            _logger.LogInformation("Finish Producing News");
        }
    }
}