namespace DevNews.Cli.Articles

open DevNews.Core.HackerNews
open DevNews.Core.Service
open DevNews.Core.UseCase
open Discord.Webhook
open MongoDB.Driver

module CompositionRoot =
    open DevNews.Core
    open DevNews.Infrastructure.Persistence
    open DevNews.Infrastructure.Notifications
    
    type T = { ParseArticlesAndNotify: ParseArticlesAndNotify }
    
    let private mongoClient (connectionString: string) =
        new MongoClient(connectionString)
    
    let private discordClient(url: string) = new DiscordWebhookClient(url)
    
    let private parseArticles = Services.parse
    
    let private providerArticles(parse: ParseArticles)= Service.provideNewArticles(parse)
    
    let private insertMany client = Repository.insertArticles client
 
    let private checkExistence client = Repository.checkArticleExistence client
    
    let private notify (discordClient: DiscordWebhookClient) = DiscordWebHooks.notify discordClient
    
    let private parseArticleAndNotify (client: MongoClient) (discordClient: DiscordWebhookClient) =
        UseCase.parseArticlesAndNotify (providerArticles(parseArticles)(checkExistence(client))) (insertMany client) (notify discordClient)
        
    let create(mongoConnectionString: string, discordWebHookUrl: string) =
        let client = mongoClient mongoConnectionString
        let discord = discordClient discordWebHookUrl
        { ParseArticlesAndNotify = parseArticleAndNotify(client)(discord) }