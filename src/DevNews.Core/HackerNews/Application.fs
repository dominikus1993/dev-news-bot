namespace DevNews.Core.HackerNews
open System.Threading.Tasks


module Repositories =
    open System.Collections.Generic
    open System.Threading.Tasks
    open DevNews.Core.Model
   
    type IHackerNewsRepository =
        abstract member Exists: articles: Article seq -> IAsyncEnumerable<struct (Article * ArticleExistence)>
        abstract member InsertMany: articles: Article seq -> Task


module Services =
    open FSharp.Control
    open FSharp.Control
    open HtmlAgilityPack
    open System.Linq

    [<Literal>]
    let private HackerNewsUrl = "https://news.ycombinator.com/"

    type ArticleDto = { Title: string; Link: string }

    type ParseHackerNewsArticles = unit -> IAsyncEnumerable<ArticleDto>

    type Notify = ArticleDto seq -> Task

    let private parseHtml() =
        async {
            let html = HtmlWeb()
            let! document = html.LoadFromWebAsync(HackerNewsUrl) |> Async.AwaitTask
            return query {
                for node in document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]") do
                    select {| link = node.GetAttributeValue("href", null); title = node.InnerText |}
            } 
        } 

    let parse () =
        asyncSeq {
            let! nodes = parseHtml()
            for node in nodes do
                yield { Title = node.title; Link = node.link }
        } 
        

module UseCases = 
    open FSharp.Control

    let parseArticlesAndNotify(parse: Services.ParseHackerNewsArticles) (repo: Repositories.IHackerNewsRepository) (notifier) = 
        parse()



