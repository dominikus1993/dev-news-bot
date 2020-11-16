namespace DevNews.Infrastructure.Persistence.HackerNews

open DevNews.Core.HackerNews.Repositories
open DevNews.Infrastructure.Persistence.HackerNews.Repositories
open Microsoft.Extensions.DependencyInjection

module internal IoC =
    open System
    
    let private addTransient<'a when 'a: not struct> (provider: IServiceProvider -> 'a) (services: IServiceCollection) =
        services.AddTransient<'a>(Func<IServiceProvider, 'a>(provider))

    let private addSingleton<'a when 'a: not struct> (provider: IServiceProvider -> 'a) (services: IServiceCollection) =
        services.AddSingleton<'a>(Func<IServiceProvider, 'a>(provider))
        
    let internal addHackerNews (services: IServiceCollection) =
        services.AddTransient<IHackerNewsRepository, MongoHackerNewsRepository>()