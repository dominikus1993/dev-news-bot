using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using DevNews.Core.Abstractions;
using DevNews.Core.Extensions;
using DevNews.Core.Model;
using DevNews.Core.Repository;
using Open.ChannelExtensions;

namespace DevNews.Core.Providers
{
    public class ChannelArticlesProvider : IArticlesProvider
    {
        private readonly IEnumerable<IArticlesParser> _articlesParsers;
        private readonly IArticlesRepository _articlesRepository;

        public ChannelArticlesProvider(IEnumerable<IArticlesParser> articlesParsers, IArticlesRepository articlesRepository)
        {
            _articlesParsers = articlesParsers;
            _articlesRepository = articlesRepository;
        }

        public IAsyncEnumerable<Article> Provide(CancellationToken cancellationToken = default)
        {
            var reader = StartProducing(_articlesParsers, cancellationToken);
            return reader.ReadAllAsync(cancellationToken)
                .Where(static article => article.IsValidArticle())
                .Select(static article => article.WithTrimmedTitle())
                .WhereAwait(async article => await NotExists(_articlesRepository, article, cancellationToken));
        }

        private static async Task<bool> NotExists(IArticlesRepository repository, Article article, CancellationToken cancellationToken)
        {
            return !await repository.Exists(article, cancellationToken);
        }

        private static ChannelReader<Article> StartProducing(IEnumerable<IArticlesParser> articlesParsers, CancellationToken cancellationToken)
        {
            var channel = Channel.CreateUnbounded<Article>(new UnboundedChannelOptions() { SingleReader = true, SingleWriter = false });
            var parsers = articlesParsers.Select(parser => Produce(parser, channel.Writer, cancellationToken));

            Task.Run(async () =>
            {
                try
                {
                    await Task.WhenAll(parsers);
                    channel.Writer.Complete();
                }
                catch (System.Exception ex)
                {

                    channel.Writer.Complete(ex);
                }

            }, cancellationToken);
            return channel.Reader;
        }
        
        private static async Task Produce(IArticlesParser parser, ChannelWriter<Article> writer, CancellationToken cancellationToken)
        {
            await foreach (var article in parser.Parse(cancellationToken))
            {
                await writer.WriteAsync(article, cancellationToken);
            }
        }
        
    }
}