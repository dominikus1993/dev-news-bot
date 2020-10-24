namespace DevNews.Application.HackerNews.Dto
{
    public class ArticleDto
    {
        public ArticleDto(string title, string link)
        {
            Title = title;
            Link = link;
        }

        public string Title { get; }
        public string Link { get; }

        public override string ToString()
        {
            return $"{nameof(Title)}: {Title}, {nameof(Link)}: {Link}";
        }
    }
}