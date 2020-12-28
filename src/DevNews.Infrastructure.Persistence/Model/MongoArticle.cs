using System;
using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace DevNews.Infrastructure.Persistence.Model
{
    public class MongoArticle
    {
        [BsonId, BsonRepresentation(BsonType.String)]
        public string Title { get; init; }

        [BsonElement] public string Link { get; init; }

        [BsonElement, BsonDateTimeOptions(Kind = DateTimeKind.Utc)]
        public DateTime CrawledAt { get; init; }

        [BsonElement, BsonDateTimeOptions(Kind = DateTimeKind.Utc)]
        public DateTime? PublishedAt { get; init; }
    }
}