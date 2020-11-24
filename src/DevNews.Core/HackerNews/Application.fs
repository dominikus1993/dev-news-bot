namespace DevNews.Core.HackerNews

open System
open DevNews.Core.Model
open FSharp.Control


module Repositories =

    type InsertManyResult = Async<Result<Article seq, ApplicationError>>

    type GetIfNotExists = Article -> Async<Option<Article>>

    type InsertMany = Article seq -> InsertManyResult

    type IHackerNewsRepository =
        abstract Exists: GetIfNotExists
        abstract InsertMany: InsertMany

module Services =
    open HtmlAgilityPack
    open Repositories

    [<Literal>]
    let private HackerNewsUrl = "https://news.ycombinator.com/"

    type ParseHackerNewsArticles = unit -> AsyncSeq<Article>

    type Notify = Article seq -> Async<unit>

    type Broadcast = Article seq -> Async<unit>

    type GetNewArticles = unit -> Async<Option<Article seq>>

    type IArticleParser =
        abstract Parse: ParseHackerNewsArticles

    type INotifier =
        abstract Notify: Notify

    type INotificationBroadcaster =
        abstract Broadcast: Article seq -> Async<unit>

    type INewArticlesProvider =
        abstract ProvideNewArticles: GetNewArticles
        
    let private parse () =
        asyncSeq {
            let html = HtmlWeb()

            let! document =
                html.LoadFromWebAsync(HackerNewsUrl)
                |> Async.AwaitTask

            let nodes =
                query {
                    for node in document.DocumentNode.SelectNodes("//*[@class=\"storylink\"]") do
                        select
                            struct {| link = node.GetAttributeValue("href", null)
                                      title = node.InnerText |}
                }

            for node in nodes do
                yield { Title = node.title; Link = node.link; Source = "HackerNews" }
        }

    let private getNewArticles (parse: ParseHackerNewsArticles) (getIfNotExists: GetIfNotExists) () =
        async {
            let result = parse ()
                            |> AsyncSeq.mapAsyncParallel (getIfNotExists)
                            |> AsyncSeq.choose (id)
                            |> AsyncSeq.toArrayAsync
            match! result with
            | [||] -> return None
            | articles -> return Some(articles |> Array.toSeq)
        }

    let private consoleNotifier (articles: Article seq) =
        async {
            printfn "%A" articles
            return ()
        }

    type ConsoleNotifier() =
        interface INotifier with
            member this.Notify = consoleNotifier

    type Broadcaster(notifiers: INotifier seq) =
        interface INotificationBroadcaster with
            member this.Broadcast(articles: Article seq) =
                async {
                    do! notifiers
                        |> Seq.map (fun notifier -> notifier.Notify(articles))
                        |> Seq.toArray
                        |> Async.Parallel
                        |> Async.Ignore
                }

    type HtmlArticleParser() =
        interface IArticleParser with
            member this.Parse = parse

    type HtmlNewArticlesProvider(repo: IHackerNewsRepository, parser: IArticleParser) =
        interface INewArticlesProvider with
            member this.ProvideNewArticles =
                getNewArticles (parser.Parse) (repo.Exists)

module UseCases =
    open Repositories
    open Services
    open DevNews.Utils

    type ParseHackerNewsArticlesAndNotify = unit -> Async<unit>
    type CheckPossibilityOfParsingArticles = DateTime -> Async<bool>

    let private insert (insertMany: InsertMany) (articles: Article seq) = articles |> insertMany

    let private notifyUsers (notify: Broadcast) (insertRes: InsertManyResult) = insertRes |> AsyncResult.map (notify)

    let private insertAndNotifyUser (insertDb: InsertMany) (notify: Broadcast) (articles: Article seq) =
        articles
            |> insert (insertDb)
            |> notifyUsers (notify)
            |> Async.Ignore

    let private parseArticlesAndNotify (getNewArticles: GetNewArticles) (insertMany: InsertMany) (notify: Broadcast) () =
        getNewArticles()
        |> AsyncOption.ifSome (insertAndNotifyUser insertMany notify)
    
   
    type GetNewArticlesAndNotifyUseCase(provider: INewArticlesProvider,
                                        repo: IHackerNewsRepository,
                                        notifier: INotificationBroadcaster) =
        member this.Execute =
            parseArticlesAndNotify
                (provider.ProvideNewArticles)
                (repo.InsertMany)
                (notifier.Broadcast)
