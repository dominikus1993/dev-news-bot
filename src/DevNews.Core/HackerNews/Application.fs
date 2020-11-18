namespace DevNews.Core.HackerNews

open System.Threading.Tasks
open DevNews.Core.Model


module Repositories =
    open System.Collections.Generic
    open System.Threading.Tasks
    open DevNews.Core.Model

    type CheckExistence = Article -> Async<Option<Article>>

    type InsertMany = Article seq -> Async<unit>

    type IHackerNewsRepository =
        abstract Exists: CheckExistence
        abstract InsertMany: InsertMany

    let private fakeExists (article: Article) = async { return Some(article) }

    let private fakeInsertMany (articles: Article seq) = async { return () }

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

    type Broadcast = Article seq -> Async<unit>

    type GetNewArticles = unit -> AsyncSeq<Article>

    type IArticleParser =
        abstract Parse: ParseHackerNewsArticles

    type INotifier =
        abstract Notify: Notify

    type INotificationBroadcaster =
        abstract Broadcast: Article seq -> Async<unit>

    type INewArticlesProvider =
        abstract Provide: GetNewArticles

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
                yield { Title = node.title; Link = node.link }
        }

    let private getNewArticles (parse: ParseHackerNewsArticles) (getIfExists: CheckExistence) () =
        parse ()
        |> AsyncSeq.map (fun x -> { Title = x.Title; Link = x.Link })
        |> AsyncSeq.mapAsyncParallel (fun x -> getIfExists (x))
        |> AsyncSeq.choose (id)

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
            member this.Provide =
                getNewArticles (parser.Parse) (repo.Exists)

module UseCases =
    open FSharp.Control
    open Repositories
    open Services
    open DevNews.Utils

    type ParseHackerNewsArticlesAndNotify = unit -> Async<unit>

    let private parseArticles (getNewArticles: GetNewArticles) =
        async {
            match! getNewArticles () |> AsyncSeq.toArrayAsync with
            | [||] -> return None
            | articles -> return Some(articles)
        }

    let private insertAndNotifyUser (insertMany: InsertMany) (notify: Broadcast) (articles: Article array) =
        async {
            do! insertMany articles
            do! notify articles
        }

    let private parseArticlesAndNotify (getNewArticles: GetNewArticles) (insertMany: InsertMany) (notify: Broadcast) () =
        parseArticles (getNewArticles)
            |> AsyncOption.ifSome (insertAndNotifyUser insertMany notify)

    type GetNewArticlesAndNotifyUseCase(provider: INewArticlesProvider,
                                        repo: IHackerNewsRepository,
                                        notifier: INotificationBroadcaster) =
        member this.Execute =
            parseArticlesAndNotify (provider.Provide) (repo.InsertMany) (notifier.Broadcast)
