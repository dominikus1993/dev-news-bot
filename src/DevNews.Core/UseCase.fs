namespace DevNews.Core

open DevNews.Utils
open FSharp.Control

module UseCase =
    open System
    open DevNews.Core.Model
    open DevNews.Core.Service
    open DevNews.Core.Repository
    
    type ParseHackerNewsArticlesAndNotify = unit -> Async<Result<unit, ApplicationError>>
    type CheckPossibilityOfParsingArticles = DateTime -> Async<bool>

    let private getNewArticles (provideNewArticles: ProvideNewArticles)() =
        async {
            let result = provideNewArticles () |> AsyncSeq.toArrayAsync
            match! result with
            | [||] -> return None
            | articles -> return Some(articles |> Array.toSeq)
        }
        
    let private insert (insertMany: InsertMany) (articles: Article seq) = articles |> insertMany

    let private notifyUsers (notify: Notify) (insertRes: InsertManyResult) = insertRes |> AsyncResult.map (notify)
    
    let private insertAndNotifyUser (insertDb: InsertMany) (notify: Notify) (articles: Article seq) =
        articles
            |> insert (insertDb)
            |> notifyUsers (notify)
            |> Async.Ignore
            
    let parseArticlesAndNotify(provideNewArticles: ProvideNewArticles)(insertMany: InsertMany)(notify: Notify) =
        let provider = getNewArticles(provideNewArticles)
        provider()
            |> AsyncOption.ifSome (insertAndNotifyUser insertMany notify)
        
