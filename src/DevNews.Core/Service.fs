namespace DevNews.Core

open System
open DevNews.Core.Repository
open DevNews.Utils

module Service =
    open DevNews.Core.Model
    open FSharp.Control
    type ParseArticles = unit -> AsyncSeq<Article>
    type ProvideNewArticles = unit -> AsyncSeq<Article>
    type Notify = Article seq -> Async<unit>

    let private getArticleIfNotExists (checkArticleExistence: CheckArticleExistence) (article: Article) =
        async {
            let! exists = checkArticleExistence (article)
            return if exists then None else Some(article)
        }
    
    let private sortBy by sequence =
        asyncSeq {
            let! res = sequence |> AsyncSeq.toArrayAsync
            for r in res |> Seq.sortBy(by) do
                yield r
        }
    let provideNewArticles (parseArticles: ParseArticles)
                           (checkArticleExistence: CheckArticleExistence)
                           ()
                           =
        parseArticles ()
                |> AsyncSeq.mapAsyncParallel (getArticleIfNotExists (checkArticleExistence))
                |> AsyncSeq.choose (id)
        
