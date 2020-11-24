namespace DevNews.Core

open FSharp.Control

module UseCase =
    open System
    open DevNews.Core.Service
    
    type ParseHackerNewsArticlesAndNotify = unit -> Async<unit>
    type CheckPossibilityOfParsingArticles = DateTime -> Async<bool>

    let private getNewArticles (provideNewArticles: ProvideNewArticles)() =
        async {
            let result = provideNewArticles () |> AsyncSeq.toArrayAsync
            match! result with
            | [||] -> return None
            | articles -> return Some(articles |> Array.toSeq)
        }
    let parseArticlesAndNotify(provideNewArticles: ProvideNewArticles)(notify: Notify) =
        async {
            return ()
        }
        
