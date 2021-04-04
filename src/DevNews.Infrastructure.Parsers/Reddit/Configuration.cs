namespace DevNews.Infrastructure.Parsers.Reddit
{
    internal class RedditConfiguration
    {
        public string[]? SubReddits { get; set; }
        public string AppId { get; set; } = "";
        public string Secret { get; set; } = "";
    }
}