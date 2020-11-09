namespace DevNews.Infrastructure.Notifications

module Discord =
    open DevNews.Core.Model
    open DevNews.Core.HackerNews.Services
    
    let private notify(articles: Article seq)  =
        async {
                      
        }
    
    type WebHookNotifier() =
        interface INotifier with
            member this.Notify =

