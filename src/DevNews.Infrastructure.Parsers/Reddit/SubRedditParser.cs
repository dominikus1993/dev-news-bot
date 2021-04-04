using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Model;
using Reddit;

namespace DevNews.Infrastructure.Parsers.Reddit
{
    internal sealed class SubRedditParser
    {
        private const int PostToDownload = 10;
        private readonly RedditClient _client;

        public SubRedditParser(RedditClient client)
        {
            _client = client;
        }

        public ValueTask<Article[]> Parse(string name)
        {
            return new(_client.Subreddit(name).Posts.New.Take(PostToDownload)
                .Select(post => new Article(post.Title, post.Permalink))
                .ToArray());
        }
    }
}