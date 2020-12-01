using DevNews.Core.Model;
using DevNews.Infrastructure.Notifications.Discord;
using FluentAssertions;
using Xunit;

namespace DevNews.Infrastructure.Notifications.UnitTests.Discord
{
    public class ExtensionsTests
    {
        [Fact]
        public void CreateEmbedWhenContentIsNull()
        {
            var article = new Article("test", "https://dotnetomaniak.pl/Gitlab-CI-Paczki-Nuget-bd90");
            var subject = article.CreateEmbed();
            subject.Description.Should().BeNullOrEmpty();
            subject.Url.Should().NotBeNullOrEmpty();
            subject.Title.Should().NotBeNullOrEmpty();
            subject.Url.Should().Be("https://dotnetomaniak.pl/Gitlab-CI-Paczki-Nuget-bd90");
            subject.Title.Should().Be("test");
        }
        
        [Fact]
        public void CreateEmbedWhenContentIsNotNull()
        {
            var article = new Article("test", "test", "https://dotnetomaniak.pl/Gitlab-CI-Paczki-Nuget-bd90");
            var subject = article.CreateEmbed();
            subject.Url.Should().NotBeNullOrEmpty();
            subject.Title.Should().NotBeNullOrEmpty();
            subject.Description.Should().NotBeNullOrEmpty();
            subject.Url.Should().Be("https://dotnetomaniak.pl/Gitlab-CI-Paczki-Nuget-bd90");
            subject.Title.Should().Be("test");
            subject.Description.Should().Be("test");
        }
    }
}