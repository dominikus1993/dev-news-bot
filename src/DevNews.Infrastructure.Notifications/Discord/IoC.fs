namespace DevNews.Infrastructure.Notifications.Discord

open System
open DevNews.Core.HackerNews.Services
open DevNews.Infrastructure.Notifications.DiscordWebHooks
open Discord.Webhook
open Microsoft.Extensions.Configuration
open Microsoft.Extensions.DependencyInjection

module IoC =
    let private addSingleton<'a when 'a: not struct> (provider: IServiceProvider -> 'a) (services: IServiceCollection) =
        services.AddSingleton<'a>(Func<IServiceProvider, 'a>(provider))
       
        
    let addDiscord (services: IServiceCollection) =
        services
            |> addSingleton(fun sp -> new DiscordWebhookClient(sp.GetService<IConfiguration>().GetConnectionString("Discord")))
            |> ignore
            
        services.AddScoped<INotifier, DiscordWebHookNotifier>()
