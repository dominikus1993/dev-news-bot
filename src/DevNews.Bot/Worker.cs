using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Akka.Actor;
using DevNews.Application.Framework.Akka;
using DevNews.Application.HackerNews.Actors;
using DevNews.Application.HackerNews.UseCases;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace DevNews.DiscordBot
{
    public class Worker : BackgroundService
    {
        private readonly ILogger<Worker> _logger;
        private readonly ParseHackerNewsMainPageAndNotifyUsersUseCase _parseHackerNewsMainPageAndNotifyUsersUseCase;

        public Worker(ILogger<Worker> logger, ParseHackerNewsMainPageAndNotifyUsersUseCase parseHackerNewsMainPageAndNotifyUsersUseCase)
        {
            _logger = logger;
            _parseHackerNewsMainPageAndNotifyUsersUseCase = parseHackerNewsMainPageAndNotifyUsersUseCase;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                await _parseHackerNewsMainPageAndNotifyUsersUseCase.Execute();
                _logger.LogInformation("Message Send");
                await Task.Delay(10000, stoppingToken);
            }
        }
    }
}
