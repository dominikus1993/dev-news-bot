using System.Threading.Tasks;
using Cocona;
using Cocona.Application;
using DevNews.Cli.Infrastructure;
using DevNews.Core.DependencyInjection;
using DevNews.Core.UseCases;
using DevNews.Infrastructure.Notifications.DependencyInjection;
using DevNews.Infrastructure.Parsers.DependencyInjection;
using DevNews.Infrastructure.Persistence.DependencyInjection;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;

namespace DevNews.Cli
{
    public class Program
    {
        public static async Task Main(string[] args) =>
            await CoconaApp.Create()
                .UseLogger("DevNews.Cli")
                .ConfigureServices((ctx, services) =>
                {
                    var configuration = ctx.Configuration;
                    services.AddPersistence(configuration);
                    services.AddNotifiers(configuration);
                    services.AddParsers(configuration);
                    services.AddCore();
                    services.AddTransient<ProduceNewsCommand>();
                })
                .ConfigureAppConfiguration(builder =>
                {
                    builder
                        .AddJsonFile("appsettings.json", true, true)
                        .AddEnvironmentVariables();
                })
                .RunAsync<Program>(args);

        public async Task
            ProduceNews(int articleQuantity, [FromService] ProduceNewsCommand command)
        {
            await command.ProduceNews(articleQuantity);
        }
    }


    public sealed class ProduceNewsCommand
    {
        private readonly ILogger<Program> _logger;
        private readonly ICoconaAppContextAccessor _coconaAppContextAccessor;
        private readonly ParseArticlesAndSendItUseCase _articlesAndSendItUseCase;

        public ProduceNewsCommand(ILogger<Program> logger, ICoconaAppContextAccessor coconaAppContextAccessor, ParseArticlesAndSendItUseCase articlesAndSendItUseCase)
        {
            _logger = logger;
            _coconaAppContextAccessor = coconaAppContextAccessor;
            _articlesAndSendItUseCase = articlesAndSendItUseCase;
        }

        public async Task ProduceNews(int articleQuantity)
        {
            _logger.LogInformation("Start Producing News");
            await _articlesAndSendItUseCase.Execute(new ParseArticlesAndSendItParam(articleQuantity), _coconaAppContextAccessor.Current.CancellationToken);
            _logger.LogInformation("Finish Producing News");
        }
    }
}