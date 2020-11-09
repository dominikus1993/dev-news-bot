namespace DevNews.Core.HackerNews

open DevNews.Core.HackerNews.Repositories
open DevNews.Core.HackerNews.Services
open DevNews.Core.HackerNews.UseCases
open Microsoft.Extensions.DependencyInjection
open System

module IoC =
    let private addTransient<'a when 'a: not struct> (provider: IServiceProvider -> 'a) (services: IServiceCollection) =
        services.AddTransient<'a>(Func<IServiceProvider, 'a>(provider))

    let addHackerNews (services: IServiceCollection) =

        services.AddScoped<IHackerNewsRepository, FakeHackerNewsRepository>() |> ignore
        services.AddScoped<IArticleParser, HtmlArticleParser>() |> ignore
        services.AddScoped<INotifier, ConsoleNotifier>() |> ignore
        services.AddScoped<INewArticlesProvider, HtmlNewArticlesProvider>() |> ignore
        services.AddScoped<GetNewArticlesAndNotifyUseCase>() |> ignore
        services
