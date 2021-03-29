using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;

namespace DevNews.Infrastructure.Parsers.Reddit
{
    internal class RedditParser : IArticlesParser
    {
        private RedditConfiguration _redditConfiguration;
        private SubRedditParser _subRedditParser;

        public RedditParser(RedditConfiguration redditConfiguration, SubRedditParser subRedditParser)
        {
            _redditConfiguration = redditConfiguration;
            _subRedditParser = subRedditParser;
        }

        public async IAsyncEnumerable<Article> Parse()
        {
            var subs = _redditConfiguration.SubReddits;
            if (subs is null)
            {
                throw new AggregateException(nameof(_redditConfiguration.SubReddits));
            }

            var articles = await Task.WhenAll(subs.Select(sub => _subRedditParser.Parse(sub)));

            foreach (var article in articles.SelectMany(x => x))
            {
                yield return article;
            }
        }
    }
}