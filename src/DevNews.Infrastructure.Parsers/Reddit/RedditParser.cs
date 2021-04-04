using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using LanguageExt;
using Open.ChannelExtensions;

[assembly: InternalsVisibleTo("DevNews.Infrastructure.Parsers.UnitTests")]
namespace DevNews.Infrastructure.Parsers.Reddit
{
    internal class RedditParser : IArticlesParser
    {
        private readonly RedditConfiguration _redditConfiguration;
        private readonly SubRedditParser _subRedditParser;

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

            var subRedditPostsChannel = subs.ToChannel()
                .TaskPipeAsync(Environment.ProcessorCount, async sub => await _subRedditParser.Parse(sub))
                .Pipe(static articlesOption => articlesOption.IfNone(static () => Array.Empty<Article>()))
                .Join();

            await foreach (var article in subRedditPostsChannel.ReadAllAsync())
            {
                yield return article;
            }
        }
    }
}