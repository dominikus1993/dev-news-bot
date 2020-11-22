namespace DevNews.Core

module UseCase =
    open System
    
    type ParseHackerNewsArticlesAndNotify = unit -> Async<unit>
    type CheckPossibilityOfParsingArticles = DateTime -> Async<bool>
    
    let parseArticlesAndNotify (articleProviders: I) =
        
