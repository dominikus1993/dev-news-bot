using System;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;
using Serilog;
using Serilog.Events;
using Serilog.Exceptions;
using Serilog.Formatting.Compact;
using Serilog.Formatting.Elasticsearch;
using Serilog.Sinks.SystemConsole.Themes;

namespace DevNews.DiscordBot.Infrastructure.Serilog
{
    public static class LoggingExtensions
    {
        public static IHostBuilder UseLogger(this IHostBuilder hostBuilder, string? applicationName = null)
        {
            return hostBuilder.UseSerilog(((context, configuration) =>
            {
                var serilogOptions = context.Configuration.GetSection("Serilog").Get<SerilogOptions>() ?? new SerilogOptions();
                if (!Enum.TryParse<LogEventLevel>(serilogOptions.MinimumLevel, true, out var level))
                {
                    level = LogEventLevel.Information;
                }

                var conf = configuration
                    .MinimumLevel.Is(level)
                    .Enrich.FromLogContext()
                    .Enrich.WithProperty("Environment", context.HostingEnvironment.EnvironmentName)
                    .Enrich.WithProperty("ApplicationName", applicationName)
                    .Enrich.WithEnvironmentUserName()
                    .Enrich.WithProcessId()
                    .Enrich.WithProcessName()
                    .Enrich.WithThreadId()
                    .Enrich.WithExceptionDetails();

                conf.WriteTo.Async((logger) =>
                {

                    if (serilogOptions.ConsoleEnabled)
                    {
                        switch (serilogOptions.Format.ToLower())
                        {
                            case "elasticsearch":
                                logger.Console(new ElasticsearchJsonFormatter());
                                break;
                            case "compact":
                                logger.Console(new CompactJsonFormatter());
                                break;
                            case "colored":
                                logger.Console(theme: AnsiConsoleTheme.Code);
                                break;
                        }
                    }

                    logger.Trace();
                });
            }));
        }
    }
}