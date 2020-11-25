namespace DevNews.Core.HackerNews

module Services =
    open HtmlAgilityPack
    open System
    open DevNews.Core.Model
    open FSharp.Control
    
    [<Literal>]
    let private HackerNewsUrl = "https://news.ycombinator.com/"
        
    let parse () =
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

