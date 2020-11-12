namespace DevNews.Infrastructure.Persistence.HackerNews
open MongoDB.Bson.Serialization.Attributes
open MongoDB.Bson

[<CLIMutable>]
type MongoArticle = { [<BsonId; BsonRepresentation(BsonType.String)>] Title: string; [<BsonElement>] Link: string }
    

