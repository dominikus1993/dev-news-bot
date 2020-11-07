namespace DevNews.Core.HackerNews
open System.Threading.Tasks
open DevNews.Core.Model


module Repositories =
    open System.Collections.Generic
    open System.Threading.Tasks
    open DevNews.Core.Model
    
    type Exists = Article -> Async<Option<Article>>
    
    type InsertMany = Article seq -> Async<unit>


module Services =
    open FSharp.Control
    open FSharp.Control
    open HtmlAgilityPack
    open System.Linq

    [<Literal>]
    let private HackerNewsUrl = "https://news.ycombinator.com/"

    type ParseHackerNewsArticles = unit -> AsyncSeq<Article>

    type Notify = Article seq -> Async<unit>
    
    type GetNewArticles  = unit -> AsyncSeq<Article>

    let parse () =
        asyncSeq {
            let html = HtmlWeb()
            let! document = html.LoadFromWebAsync(HackerNewsUrl) |> Async.AwaitTask
            let nodes = query {
                for node in document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]") do
                    select struct {| link = node.GetAttributeValue("href", null); title = node.InnerText |}
            }
            
            for node in nodes do
                yield { Title = node.title; Link = node.link }
        } 
        
    let getNewArticles (parse: ParseHackerNewsArticles) (getIfExists: Repositories.Exists) =
        parse()
            |> AsyncSeq.map(fun x -> { Title = x.Title; Link = x.Link })
            |> AsyncSeq.mapAsyncParallel(fun x -> getIfExists(x))
            |> AsyncSeq.choose(id)
   
module UseCases = 
    open FSharp.Control
    
    
    
    let parseArticlesAndNotify(getNewArticles: Services.GetNewArticles) (insertMany: Repositories.InsertMany)(notify: Services.Notify) =
        async {
            let! articles = getNewArticles() |> AsyncSeq.toArrayAsync
            do! insertMany articles
            do! notify articles
            return ();
        }



