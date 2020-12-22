using DevNews.Core.Model;
using FluentAssertions;
using Xunit;

namespace DevNews.Core.UnitTests.Model
{
    public class ArticleRecordTests
    {
        [Fact]
        public void TestArticleValidationWhenUrlIsInCorrect()
        {
            var article = new Article("test", "notlink");
            var subject = article.IsValidArticle();

            subject.Should().BeFalse();
        }
        
        [Fact]
        public void TestArticleValidationWhenUrlIsCorrect()
        {
            var article = new Article("test", "https://scienceintegritydigest.com/2020/12/20/paper-about-herbalife-related-patient-death-removed-after-company-threatens-to-sue-the-journal/");
            var subject = article.IsValidArticle();

            subject.Should().BeTrue();
        }
    }
}