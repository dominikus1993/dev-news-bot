namespace DevNews.Bot

open Microsoft.AspNetCore.Builder

module App =
    
    open Microsoft.Extensions.Hosting
    open Microsoft.Extensions.Logging
    open FSharp.Control.Tasks.V2.ContextInsensitive
    open System.Threading
    open Microsoft.Extensions.DependencyInjection
    open Saturn
    open DevNews.Core.HackerNews
    
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
        services.AddHostedService<HackerNewsWorker>() |> ignore
        services

//    let configureApp(app: IApplicationBuilder) =
//        app.Add
    [<EntryPoint>]
    let main argv =
        let h =
            application {
                no_webhost //Don't start default webhost
                cli_arguments argv
                service_config configureServices
            }
        run h
        0 // return an integer exit code