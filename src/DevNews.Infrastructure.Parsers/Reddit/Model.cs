using System.Text.Json.Serialization;

namespace DevNews.Infrastructure.Parsers.Reddit
{
    internal class Data
    {
        [JsonPropertyName("children")] 
        public PostData[]? Posts { get; init; }
    }

    internal record PostData
    {
        [JsonPropertyName("data")] 
        public Post? Post { get; init; }
    }

    public class Post
    {
        [JsonPropertyName("title")] 
        public string? Title { get; init; }   
        [JsonPropertyName("url")] 
        public string? Url { get; init; } 
        [JsonPropertyName("selftext")] 
        public string? Content { get; set; }
    }
    internal class Subreddit
    {
        [JsonPropertyName("data")]
        public Data? Data { get; init; }
    };
}