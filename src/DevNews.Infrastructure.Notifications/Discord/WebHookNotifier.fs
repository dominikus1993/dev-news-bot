namespace DevNews.Infrastructure.Notifications

module DiscordWebHooks =
    open Discord.Webhook
    open Discord
    open DevNews.Core.Model
    open DevNews.Core.HackerNews.Services
    open System
    let private options = RequestOptions(Timeout = Nullable(15000))
        
    let private notify(client: DiscordWebhookClient)(articles: Article seq) =
        async {
            let embeds = articles
                            |> Seq.map(fun article -> EmbedBuilder().WithUrl(article.Link).WithTitle(article.Title))
                            |> Seq.map(fun x -> x.Build())
                            |> Seq.toArray
            let! res = client.SendMessageAsync("Nowe newsy od HackerNews", false, embeds, options = options) |> Async.AwaitTask 
            return ()
        }
    
    type DiscordWebHookNotifier(client: DiscordWebhookClient) =
        interface INotifier with
            member this.Notify = notify(client)

