namespace DevNews.Bot

open System
open DevNews.Infrastructure.Persistence
open Microsoft.AspNetCore.Builder

module App =
    
    open Microsoft.Extensions.Hosting
    open Microsoft.Extensions.Logging
    open FSharp.Control.Tasks.V2.ContextInsensitive
    open System.Threading
    open Microsoft.Extensions.DependencyInjection
    open Saturn
    open DevNews.Core.HackerNews
    open Microsoft.Extensions.Configuration
    
    let configureServices (services: IServiceCollection) =
        services

    let configureHost(host: IHostBuilder) =
        host.ConfigureAppConfiguration(fun builder -> builder.AddJsonFile("appsettings.json", true, true).AddUserSecrets("caad687c-126f-440c-8dc2-c85d8ad6668a").AddEnvironmentVariables() |> ignore)
        
    [<EntryPoint>]
    let main argv =
        let h =
            application {
                no_webhost //Don't start default webhost
                cli_arguments argv
                service_config configureServices
                host_config configureHost
            }
        run h
        0 // return an integer exit code