namespace DevNews.Core

open DevNews.Core.Repository
open DevNews.Utils

module Service =
    open DevNews.Core.Model
    open FSharp.Control
    
    type ProvideArticles = unit -> AsyncSeq<Article>

    type Notify = Article seq -> Async<unit>

    type GetNewArticles = unit -> Async<Option<Article seq>>

