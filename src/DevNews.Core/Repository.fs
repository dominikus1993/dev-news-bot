namespace DevNews.Core

module Repository =
    open DevNews.Core.Model
    
    type InsertManyResult = Async<Result<Article seq, ApplicationError>>

    type CheckArticleExistence = Article -> Async<bool>

    type InsertMany = Article seq -> InsertManyResult