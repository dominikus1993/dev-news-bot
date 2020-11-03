namespace DevNews.Core.HackerNews


module Repositories =
    open System.Collections.Generic
    open System.Threading.Tasks
    open DevNews.Core.Model
   
    type IHackerNewsRepository =
        abstract member Exists: articles: Article seq -> IAsyncEnumerable<struct (Article * ArticleExistence)>
        abstract member InsertMany: articles: Article seq -> Task


module Services =
    open FSharp.Control
    open HtmlAgilityPack
    open System.Linq

    [<Literal>]
    let private HackerNewsUrl = "https://news.ycombinator.com/"

    type ArticleDto = { Title: string; Link: string }

    type ParseHackerNewsArticles = unit -> IAsyncEnumerable<ArticleDto>

    let parse () =
        asyncSeq {
            let html = HtmlWeb()
            let! document = html.LoadFromWebAsync(HackerNewsUrl) |> Async.AwaitTask
            let nodes = document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]").Select(fun e -> {| link = e.GetAttributeValue("href", null); title = e.InnerText |})
            for node in nodes do
                yield { Title = node.title; Link = node.link }
        }
        

module UseCases = 
    let parseArticlesAndNotify = 
        2