namespace DevNews.Infrastructure.Persistence

module Repository =
    open DevNews.Core.Repository
    open System
    open DevNews.Core.Model
    open FSharp.Control
    open MongoDB.Driver
    open MongoDB.Driver.Linq

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
        
    let private insertArticlesToDb (col: IMongoCollection<MongoArticle>) (articles: Article seq) : InsertManyResult =
        async {
            let writeMode = articles |> Seq.map(fun x -> InsertOneModel<MongoArticle>({ Link = x.Link; Title = x.Title; CrawledAt = DateTime.UtcNow; PublishedAt = Nullable()}) :>  WriteModel<MongoArticle>)
            do! col.BulkWriteAsync(writeMode) |> Async.AwaitTask |> Async.Ignore
            return Ok(articles)
        }
        
    let private getDb(client: IMongoClient) =
        client.GetDatabase("HackerNews")
    
    let private getCollection(db: IMongoDatabase) =
        db.GetCollection<MongoArticle>("articles")
    
    let insertArticles(client: IMongoClient) =
        client |> getDb |> getCollection |> insertArticlesToDb
        
    let checkArticleExistence(client: IMongoClient) =
        client |> getDb |> getCollection |> checkExistence

