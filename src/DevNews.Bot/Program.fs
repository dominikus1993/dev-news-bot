namespace DevNews.Bot

open Microsoft.AspNetCore.Builder

module App =
    
    open Microsoft.Extensions.Hosting
    open Microsoft.Extensions.Logging
    open FSharp.Control.Tasks.V2.ContextInsensitive
    open System.Threading
    open Microsoft.Extensions.DependencyInjection
    open DevNews.Infrastructure.Notifications.Discord
    open Saturn
    open DevNews.Core.HackerNews
    open Microsoft.Extensions.Configuration
    
    type HackerNewsWorker(logger:ILogger<HackerNewsWorker>, useCase: UseCases.GetNewArticlesAndNotifyUseCase) =
        inherit BackgroundService()
        override __.ExecuteAsync(ct: CancellationToken) =
                ct.Register(fun () -> logger.LogInformation("Worker canceled at: {time}", System.DateTimeOffset.Now)) |> ignore
                task {
                    while not ct.IsCancellationRequested do
                    logger.LogInformation("Worker running at: {time}", System.DateTimeOffset.Now)
                    do! useCase.Execute()
                    do! Tasks.Task.Delay(10000, ct)
                } :> Tasks.Task
    
    let configureServices (services: IServiceCollection) =
        services |> IoC.addHackerNews |> ignore
        services |> IoC.addDiscord |> ignore
        services.AddHostedService<HackerNewsWorker>() |> ignore
        services.Add
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