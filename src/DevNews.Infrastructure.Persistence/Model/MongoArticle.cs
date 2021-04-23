using System;
using DevNews.Core.Model;
using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace DevNews.Infrastructure.Persistence.Model
{
    internal class MongoArticle
    {
        [BsonId, BsonRepresentation(BsonType.String)]
        public string? Title { get; init; }

        [BsonElement] public string? Link { get; init; }
        
        [BsonElement] public string? Content { get; init; }

        [BsonElement, BsonDateTimeOptions(Kind = DateTimeKind.Utc)]
        public DateTime CrawledAt { get; init; }

        [BsonElement, BsonDateTimeOptions(Kind = DateTimeKind.Utc)]
        public DateTime? PublishedAt { get; init; }

        public MongoArticle()
        {
            
        }

        public MongoArticle(Article article)
        {
            var (title, content, link) = article;
            Title = title;
            Link = link;
            Content = content;
            CrawledAt = DateTime.UtcNow;
            PublishedAt = DateTime.UtcNow;
        }

        public Article AsArticle()
        {
            return new(Title ?? "", Content, Link ?? "");
        }
    }
}