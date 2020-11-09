namespace DevNews.Core.HackerNews
open System.Threading.Tasks
open DevNews.Core.Model


module Repositories =
    open System.Collections.Generic
    open System.Threading.Tasks
    open DevNews.Core.Model
    
    type Exists = Article -> Async<Option<Article>>
    
    type InsertMany = Article seq -> Async<unit>
    
    type IHackerNewsRepository =
        abstract member Exists: Exists
        abstract member InsertMany: InsertMany

    let private fakeExists(article: Article) =
        async {
            return Some(article)
        }
        
    let private fakeInsertMany(articles: Article seq) =
        async {
            return ()
        }
    
    type FakeHackerNewsRepository() =
        interface IHackerNewsRepository with
            member this.Exists = fakeExists
            member this.InsertMany = fakeInsertMany
        
module Services =
    open FSharp.Control
    open FSharp.Control
    open HtmlAgilityPack
    open System.Linq
    open Repositories

    [<Literal>]
    let private HackerNewsUrl = "https://news.ycombinator.com/"

    type ParseHackerNewsArticles = unit -> AsyncSeq<Article>

    type Notify = Article seq -> Async<unit>
    
    type GetNewArticles  = unit -> AsyncSeq<Article>
    
    type IArticleParser =
        abstract member Parse: ParseHackerNewsArticles
        
    type INotifier =
        abstract member Notify: Notify
        
    type INewArticlesProvider =
        abstract member Provide: GetNewArticles
    let private parse () =
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
        
    let private getNewArticles (parse: ParseHackerNewsArticles) (getIfExists: Repositories.Exists)() =
        parse()
            |> AsyncSeq.map(fun x -> { Title = x.Title; Link = x.Link })
            |> AsyncSeq.mapAsyncParallel(fun x -> getIfExists(x))
            |> AsyncSeq.choose(id)
            
    let private consoleNotifier(articles: Article seq) =
        async {
            printfn "%A" articles
            return ()
        }
        
    type ConsoleNotifier() =
        interface INotifier with
            member this.Notify = consoleNotifier
            
    type HtmlArticleParser() =
        interface IArticleParser with
            member this.Parse = parse
            
    type HtmlNewArticlesProvider(repo: IHackerNewsRepository, parser: IArticleParser) =
        interface INewArticlesProvider with
            member this.Provide = getNewArticles(parser.Parse)(repo.Exists)
   
module UseCases = 
    open FSharp.Control
    open Repositories
    open Services
    
    type ParseHackerNewsArticlesAndNotify = unit -> Async<unit>
    
    let private parseArticlesAndNotify(getNewArticles: Services.GetNewArticles) (insertMany: Repositories.InsertMany)(notify: Services.Notify)()=
        async {
            let! articles = getNewArticles() |> AsyncSeq.toArrayAsync
            do! insertMany articles
            do! notify articles
            return ();
        }
        
    type GetNewArticlesAndNotifyUseCase(provider: INewArticlesProvider, repo: IHackerNewsRepository, notifier: INotifier) =
        member this.Execute = parseArticlesAndNotify(provider.Provide)(repo.InsertMany)(notifier.Notify)



