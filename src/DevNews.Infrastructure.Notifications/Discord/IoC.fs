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
            |> addSingleton(fun sp -> new DiscordWebhookClient("https://discordapp.com/api/webhooks/775805301873704991/4aKBJiLJXtQ6skHCnh0GDMdXMJo_yv7RVfNyJ9nxKANi9OFeHKyBPI_xgFdrA2oG9gBi"))
            |> ignore
            
        services.AddScoped<INotifier, DiscordWebHookNotifier>()
