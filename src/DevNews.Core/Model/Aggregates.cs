namespace DevNews.Core.Model
{
    public record Article(string Tile, string? Content, string Link)
    {
        public Article(string title, string link) : this(title, null, link)
        {
        }
    }
}