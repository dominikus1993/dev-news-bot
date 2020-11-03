namespace DevNews.Core.HackerNews

open System.Collections.Generic
open System.Threading.Tasks
open DevNews.Core.Model

type ArticleExistence = (Article * bool)

type IHackerNewsRepository =
    abstract member Exists: articles: Article seq -> IAsyncEnumerable<struct (Article * bool)>
    abstract member InsertMany: articles: Article seq -> Task