using DevNews.Core.Model;

namespace DevNews.Core.Dto
{
    public record ArticleDto(string Title, string? Content, string Link)
    {
        public ArticleDto(Article article) : this(article.Title, article.Content, article.Link)
        {
            
        }
    }
}