namespace DevNews.Cli.Articles

open MongoDB.Driver

module CompostionRoot =
    open DevNews.Core
    
    let mongoClient () = new MongoClient()
    
    let private providerArticles = HackerNews.Services.parse
    
    let parseArticleAndNotify () = UseCase.parseArticlesAndNotify (providerArticles)