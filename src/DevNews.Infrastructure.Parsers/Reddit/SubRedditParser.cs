using System;
using System.Linq;
using System.Net.Http;
using System.Net.Http.Json;
using System.Threading.Tasks;
using DevNews.Core.Model;
using LanguageExt;
using static LanguageExt.Prelude;

namespace DevNews.Infrastructure.Parsers.Reddit
{
    internal sealed class SubRedditParser
    {
        private const string PostToDownload = "10";
        private readonly HttpClient _client;

        public SubRedditParser(HttpClient client)
        {
            _client = client;
        }

        public async Task<Option<Article[]>> Parse(string name)
        {
            var url = $"r/{name}/top/.json?limit={PostToDownload}";
            var result = await _client.GetFromJsonAsync<Subreddit>(url);
            if (result?.Data?.Posts is null)
            {
                return None;
            }

            return result.Data.Posts
                .Where(x => x.Post is not null)
                .Select(x => new Article(x.Post.Title, x.Post.Content, x.Post.Url))
                .ToArray();
        }
    }
}