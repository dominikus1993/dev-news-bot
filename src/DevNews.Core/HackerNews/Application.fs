namespace DevNews.Core.HackerNews

open System.Collections.Generic
open System.Threading.Tasks
open DevNews.Core.Model

module Repositories =
    type IHackerNewsRepository =
        abstract member Exists: articles: Article seq -> IAsyncEnumerable<struct (Article * ArticleExistence)>
        abstract member InsertMany: articles: Article seq -> Task


module UseCases = 
    let parseArticlesAndNotify = 
        2