using System;
using DevNews.Core.Extensions;

namespace DevNews.Core.Model
{
    public record Article(string Title, string? Content, string Link)
    {
        public Article(string title, string link) : this(title, null, link)
        {
        }

        public Article WithTrimmedTitle() => this with {Title = Title.TrimEntersAndSpaces()};

        public bool IsValidArticle()
        {
            var contentIsValid = Content is null || Content.Length < 2048;
            var titleIsValid = !string.IsNullOrEmpty(Title);
            var linkIsValid = !string.IsNullOrEmpty(Link);
            return contentIsValid && titleIsValid && linkIsValid &&
                   Uri.TryCreate(Link, UriKind.Absolute, out var result)
                   && (result.Scheme == Uri.UriSchemeHttp || result.Scheme == Uri.UriSchemeHttps);
        }
    }
}