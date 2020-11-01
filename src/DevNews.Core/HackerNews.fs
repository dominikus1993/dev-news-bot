namespace DevNews.Core.HackerNews

open System.Collections.Generic
open DevNews.Core.Model

type ArticleExistence = (Article * bool)

type IHackerNewsRepository =
    abstract member Exists: articles: Article seq -> IAsyncEnumerable<struct (Article * bool)>