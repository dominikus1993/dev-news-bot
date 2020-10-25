using Akka.Actor;
using Akka.Configuration;
using DevNews.Application.Framework.Akka;
using DevNews.Application.HackerNews.Actors;
using DevNews.Application.HackerNews.UseCases;
using DevNews.Application.Notifications.Actors;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Application.Framework.DependencyInjection
{
    public static class IServicesCollectionExtensions
    {
        private const string Config =@"akka {
                                stream {
                                  debug-logging = on
                                  debug {
                                    fuzzing-mode = on
                                  }
                                }
                                loggers=[""Akka.Logger.Serilog.SerilogLogger, Akka.Logger.Serilog""]
                                actor {
                                  debug {
                                      receive = on
                                      autoreceive = on
                                      lifecycle = on
                                      event-stream = on
                                      unhandled = on
                                    }
                                    serializers {
                                      hyperion = ""Akka.Serialization.HyperionSerializer, Akka.Serialization.Hyperion""
                                    }
                                    serialization-bindings {
                                      ""System.Object"" = hyperion
                                    }
                                }
                                  loglevel = INFO
                                  log-config-on-start = on
                                  stdout-loglevel = INFO 
                            }";
        
        public static void AddApplication(this IServiceCollection services)
        {
            services.AddSingleton<Actors>(sp =>
            {
                var system = ActorSystem.Create("bot", ConfigurationFactory.ParseString(Config));
                var webhookActor = system.ActorOf(WebHookSenderActor.Create(sp),
                    WebHookSenderActor.HackerNewsParserActorPath.Name);
                var hackernews = system.ActorOf(HackerNewsParserActor.Create(sp),
                    HackerNewsParserActor.HackerNewsParserActorPath.Name);
                
                return new Actors(system, webhookActor, hackernews);
            });

            services.AddSingleton<ParseHackerNewsMainPageAndNotifyUsersUseCase>();

        } 
    }
}