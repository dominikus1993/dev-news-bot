namespace DevNews.Infrastructure.Persistence.HackerNews

open DevNews.Core.Model
open DevNews.Infrastructure.Persistence.HackerNews
open MongoDB.Driver
open MongoDB.Driver.Linq
module Repositories =
    open DevNews.Core.HackerNews.Repositories
    
    let checkExistence (col: IMongoCollection<MongoArticle>) (art: Article) =
        async {
            return! col.AsQueryable().AnyAsync(fun x -> x.Title = art.Title) |> Async.AwaitTask
        }