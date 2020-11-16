namespace DevNews.Infrastructure.Persistence

open Microsoft.Extensions.Configuration
open Microsoft.Extensions.DependencyInjection
open MongoDB.Driver

module IoC =
    open System
    open DevNews.Infrastructure.Persistence.HackerNews.IoC
    let private addTransient<'a when 'a: not struct> (provider: IServiceProvider -> 'a) (services: IServiceCollection) =
        services.AddTransient<'a>(Func<IServiceProvider, 'a>(provider))

    let private addSingleton<'a when 'a: not struct> (provider: IServiceProvider -> 'a) (services: IServiceCollection) =
        services.AddSingleton<'a>(Func<IServiceProvider, 'a>(provider))

    let private addMongoClient(sp: IServiceProvider) : IMongoClient =
        let config = sp.GetService<IConfiguration>()
        let connection = config.GetConnectionString("Articles")
        MongoClient(connection) :> IMongoClient
        
    let private addMongoDb (services: IServiceCollection) =
        services |> addSingleton(addMongoClient)
        
    let addPersistenceInfrastructure (services: IServiceCollection) =
        services |> addMongoDb |> ignore
        services |> addHackerNews |> ignore
        services

