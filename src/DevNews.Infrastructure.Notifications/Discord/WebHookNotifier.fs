namespace DevNews.Infrastructure.Notifications

module DiscordWebHooks =
    open Discord.Webhook
    open Discord
    open DevNews.Core.Model
    open DevNews.Core.HackerNews.Services
    open FSharp.Control.Tasks.V2
    open System
    open FSharpx.Collections
    let private options = RequestOptions(Timeout = Nullable(15000))
    
    let paged(pageSize: int)(sequence: _ seq) =
        seq {
            let mut i = 0
            use enumer = sequence.GetEnumerator()
            while enumer.MoveNext() do
                yield enumer.Current
        }
    let private notify(client: DiscordWebhookClient)(articles: Embed seq) =
        async {
            let embeds = articles
                            |> Seq.map(fun article -> EmbedBuilder().WithUrl(article.Link).WithTitle(article.Title))
                            |> Seq.map(fun x -> x.Build())
                            |> Seq.toArray
 
            do! client.SendMessageAsync("Nowe newsy od HackerNews", embeds =  embeds, options = options) |> Async.AwaitTask |> Async.Ignore
            return ()
        }
        
    let private pagedNotify(client: DiscordWebhookClient)(articles: Article seq) =
        async {
            let embeds = articles
                            |> Seq.map(fun article -> EmbedBuilder().WithUrl(article.Link).WithTitle(article.Title))
                            |> Seq.map(fun x -> x.Build())
                            |> Seq.page 101 1
                            |> Seq.toArray
                         
        }    
    type DiscordWebHookNotifier(client: DiscordWebhookClient) =
        interface INotifier with
            member this.Notify = notify(client)

