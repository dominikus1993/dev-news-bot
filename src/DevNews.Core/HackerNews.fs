namespace DevNews.Domain.HackerNews

open System.Collections.Generic
open DevNews.Domain.Model

type IHackerNewsRepository =
    abstract member Exists: articles: Article seq -> IAsyncEnumerable<struct (Article * bool)>