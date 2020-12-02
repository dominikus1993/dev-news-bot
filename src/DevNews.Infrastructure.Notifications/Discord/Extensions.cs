using DevNews.Core.Model;
using Discord;

namespace DevNews.Infrastructure.Notifications.Discord
{
    public static class Extensions
    {
        public static Embed CreateEmbed(this Article article)
        {
            var builder = article switch
            {
                var (title, content, link) when content is not null 
                    => new EmbedBuilder().WithUrl(link).WithTitle(title)
                        .WithDescription(content),
                var (title, _, link) => new EmbedBuilder().WithUrl(link).WithTitle(title)
            };
            return builder.Build();
        }
    }
}