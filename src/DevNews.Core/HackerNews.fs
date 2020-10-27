namespace DevNews.Domain.HackerNews

open System.Collections.Generic
open System.Threading.Tasks
open DevNews.Domain.Model

type IHackerNewsRepository =
    abstract member Exists: articles: Article -> ValueTask<bool>