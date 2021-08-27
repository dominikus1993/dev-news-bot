using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Threading;
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

        public async IAsyncEnumerable<Article> Parse([EnumeratorCancellation] CancellationToken cancellationToken = default)
        {
            var subs = _redditConfiguration.SubReddits;
            if (subs is null)
            {
                throw new AggregateException(nameof(_redditConfiguration.SubReddits));
            }

            var subRedditPostsChannel = subs
                .ToChannel(cancellationToken: cancellationToken)
                .TaskPipeAsync(Environment.ProcessorCount, async sub => await _subRedditParser.Parse(sub), cancellationToken: cancellationToken)
                .Pipe(static articlesOption => articlesOption.IfNone(static () => Array.Empty<Article>()), cancellationToken: cancellationToken)
                .Join();

            await foreach (var article in subRedditPostsChannel.ReadAllAsync(cancellationToken))
            {
                yield return article;
            }
        }
    }
}