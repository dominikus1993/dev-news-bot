namespace DevNews.Infrastructure.Notifications

module DiscordWebHooks =
    open Discord.Webhook
    open Discord
    open DevNews.Core.Model
    open DevNews.Core.HackerNews.Services
    open FSharp.Control.Tasks.V2
    open System
    open DevNews.Utils
    
    let private options = RequestOptions(Timeout = Nullable(15000))
        
    let private notify(client: DiscordWebhookClient)(articles: Article seq) =
        async {
            let tasks = articles
                            |> Seq.map(fun article -> EmbedBuilder().WithUrl(article.Link).WithTitle(article.Title))
                            |> Seq.map(fun x -> x.Build())
                            |> Seq.paged 10
                            |> Seq.map(fun embds -> client.SendMessageAsync("Nowe newsy od HackerNews", embeds =  embds, options = options) |> Async.AwaitTask)
                            |> Seq.toArray
            do! tasks |> Async.Parallel |> Async.Ignore             
        }    
    type DiscordWebHookNotifier(client: DiscordWebhookClient) =
        interface INotifier with
            member this.Notify = notify(client)

