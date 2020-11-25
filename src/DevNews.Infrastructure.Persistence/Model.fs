namespace DevNews.Infrastructure.Persistence

open System
open MongoDB.Bson.Serialization.Attributes
open MongoDB.Bson

[<CLIMutable>]
type private MongoArticle =
    { [<BsonId; BsonRepresentation(BsonType.String)>]
      Title: string
      [<BsonElement>]
      Link: string
      [<BsonElement; BsonDateTimeOptions(Kind = DateTimeKind.Utc)>]
      CrawledAt: DateTime
      [<BsonElement>]
      PublishedAt: Nullable<DateTime> }
