namespace DevNews.Infrastructure.Persistence.HackerNews

open DevNews.Core.Model
open MongoDB.Driver

module Repositories =
    open DevNews.Core.HackerNews.Repositories
    
    let checkExistence (mongo: IMongoClient) (article: Article) =
        async {
            let db 
        }