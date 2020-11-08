namespace DevNews.Core.HackerNews

open Microsoft.Extensions.DependencyInjection
open System

module IoC =
    let private addTransient<'a when 'a: not struct> (provider: IServiceProvider -> 'a) (services: IServiceCollection) =
        services.AddTransient<'a>(Func<IServiceProvider, 'a>(provider))

    let addHackerNews (services: IServiceCollection) =
        services
        |> addTransient<Services.ParseHackerNewsArticles> (fun _ -> Services.parse)
        |> ignore

        services
        |> addTransient<Repositories.Exists> (fun _ -> Repositories.fakeExists)
        |> ignore

        services
        |> addTransient<Repositories.InsertMany> (fun _ -> Repositories.fakeInsertMany)
        |> ignore

        services
        |> addTransient<Services.GetNewArticles> (fun sp ->
            Services.getNewArticles
                (sp.GetService<Services.ParseHackerNewsArticles>())
                (sp.GetService<Repositories.Exists>()))
        |> ignore

        services
        |> addTransient<Services.Notify> (fun _ -> Services.consoleNotifier)
        |> ignore

        services
        |> addTransient<UseCases.ParseHackerNewsArticlesAndNotify> (fun sp ->
            UseCases.parseArticlesAndNotify
                (sp.GetService<Services.ParseHackerNewsArticles>())
                (sp.GetService<Repositories.InsertMany>())
                (sp.GetService<Services.Notify>()))
        |> ignore

        services
