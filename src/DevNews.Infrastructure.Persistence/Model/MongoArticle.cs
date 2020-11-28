using System;
using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace DevNews.Infrastructure.Persistence.Model
{
    public class MongoArticle
    {
        [BsonId, BsonRepresentation(BsonType.String)]
        public string Title { get; set; }

        [BsonElement] public string Link { get; set; }

        [BsonElement, BsonDateTimeOptions(Kind = DateTimeKind.Utc)]
        public DateTime CrawledAt { get; set; }

        [BsonElement, BsonDateTimeOptions(Kind = DateTimeKind.Utc)]
        public DateTime? PublishedAt { get; set; }
    }
}