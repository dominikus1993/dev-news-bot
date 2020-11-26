namespace DevNews.Core

open DevNews.Utils
open FSharp.Control

module UseCase =
    open System
    open DevNews.Core.Model
    open DevNews.Core.Service
    open DevNews.Core.Repository

    type ParseArticlesAndNotifyParams = { ArticlesQuantity: int }

    type ParseArticlesAndNotify = ParseArticlesAndNotifyParams -> Async<unit>
    type CheckPossibilityOfParsingArticles = DateTime -> Async<bool>

    let private toOption(articles: Article array) =
        match articles with
        | [||] -> None
        | articles -> Some(articles)
        
    let filterArticles (param: ParseArticlesAndNotifyParams) (articles: Article array option) =
        articles
        |> Option.map (fun art ->
            art
            |> Seq.sortBy (fun _ -> Guid.NewGuid())
            |> Seq.truncate (param.ArticlesQuantity))

    let private getNewArticles (provideNewArticles: ProvideNewArticles) (param: ParseArticlesAndNotifyParams) =
        async {
            let! result = provideNewArticles () |> AsyncSeq.toArrayAsync
            return result |> toOption |> filterArticles(param)
        }

    let private insert (insertMany: InsertMany) (articles: Article seq) = articles |> insertMany

    let private notifyUsers (notify: Notify) (insertRes: InsertManyResult) = insertRes |> AsyncResult.map (notify)

    let private insertAndNotifyUser (insertDb: InsertMany) (notify: Notify) (articles: Article seq) =
        articles
        |> insert (insertDb)
        |> notifyUsers (notify)
        |> Async.Ignore

    let parseArticlesAndNotify (provideNewArticles: ProvideNewArticles)
                               (insertMany: InsertMany)
                               (notify: Notify)
                               (param: ParseArticlesAndNotifyParams)
                               =
        let provider = getNewArticles (provideNewArticles)

        provider (param)
        |> AsyncOption.ifSome (insertAndNotifyUser insertMany notify)
