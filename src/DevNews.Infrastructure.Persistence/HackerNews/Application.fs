namespace DevNews.Infrastructure.Persistence.HackerNews

open System
open DevNews.Core.Model
open DevNews.Infrastructure.Persistence.HackerNews
open FSharp.Control
open MongoDB.Driver
open MongoDB.Driver.Linq
module internal Repositories =
    open DevNews.Core.HackerNews.Repositories
    
    let private checkExistence (col: IMongoCollection<MongoArticle>) (art: Article) =
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
        
    let private insertArticles (col: IMongoCollection<MongoArticle>) (articles: Article seq) : InsertManyResult =
        async {
            let writeMode = articles |> Seq.map(fun x -> InsertOneModel<MongoArticle>({ Link = x.Link; Title = x.Title; CrawledAt = DateTime.UtcNow }) :>  WriteModel<MongoArticle>)
            do! col.BulkWriteAsync(writeMode) |> Async.AwaitTask |> Async.Ignore
            return Ok(articles)
        }
        
    let private getDb(client: IMongoClient) =
        client.GetDatabase("HackerNews")
    
    let private getCollection(db: IMongoDatabase) =
        db.GetCollection<MongoArticle>("articles")
        
    type MongoHackerNewsRepository(client: IMongoClient) =
        interface IHackerNewsRepository with
            member this.Exists = client |> getDb |> getCollection |> checkExistence
            
            member this.InsertMany = client |> getDb |> getCollection |> insertArticles