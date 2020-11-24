namespace DevNews.Core

module UseCase =
    open System
    open DevNews.Core.Service
    
    type ParseHackerNewsArticlesAndNotify = unit -> Async<unit>
    type CheckPossibilityOfParsingArticles = DateTime -> Async<bool>
   
        
