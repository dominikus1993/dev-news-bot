namespace DevNews.Core

open DevNews.Core.Repository
open DevNews.Utils

module Service =
    open DevNews.Core.Model
    open FSharp.Control
   
    type ParseArticles = unit -> AsyncSeq<Article>
    type ProvideNewArticles = unit -> AsyncSeq<Article>
    
    let private getArticleIfNotExists(checkArticleExistence: CheckArticleExistence)(article: Article) =
        async {
            let! exists = checkArticleExistence(article)
            return if exists then None else Some(article)
        }
    
    let provideNewArticles (parseArticles: ParseArticles)(checkArticleExistence: CheckArticleExistence)() =
        parseArticles()
            |> AsyncSeq.mapAsyncParallel(getArticleIfNotExists(checkArticleExistence))
            |> AsyncSeq.choose(id)
        