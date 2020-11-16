namespace DevNews.Infrastructure.Persistence.HackerNews

open DevNews.Core.Model
open DevNews.Infrastructure.Persistence.HackerNews
open MongoDB.Driver
open MongoDB.Driver.Linq
module Repositories =
    open DevNews.Core.HackerNews.Repositories
    
    let checkExistence (col: IMongoCollection<MongoArticle>) (art: Article) =
        async {
            let exists =
                query {
                    for article in col.AsQueryable() do
                        exists (article.Title = art.Title)
                }
            if exists then
                return None
            else
                return Some(art)
        }
        
    let insertArticles (col: IMongoCollection<MongoArticle>) (art: Article seq) =
        async {
            let writeMode = art |> Seq.map(fun x -> InsertOneModel<MongoArticle>({ Link = x.Link; Title = x.Title }) :>  WriteModel<MongoArticle>)
            do! col.BulkWriteAsync(writeMode) |> Async.AwaitTask |> Async.Ignore
        }
        
    let private getDb(client: MongoClient) =
        client.GetDatabase("HackerNews")
    
    let private getCollection(db: IMongoDatabase) =
        db.GetCollection<MongoArticle>("articles")
        
    type MongoHackerNewsRepository(client: MongoClient) =
        interface IHackerNewsRepository with
            member this.Exists = client |> getDb |> getCollection |> checkExistence
            
            member this.InsertMany = client |> getDb |> getCollection |> insertArticles