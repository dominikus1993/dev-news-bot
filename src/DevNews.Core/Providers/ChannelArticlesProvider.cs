using System.Collections.Generic;
using System.Linq;
using System.Threading.Channels;
using System.Threading.Tasks;
using DevNews.Core.Abstractions;
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

        public IAsyncEnumerable<Article> Provide()
        {
            var reader = StartProducing(_articlesParsers);
            return reader.ReadAllAsync()
                .WhereAwait(async article => await NotExists(_articlesRepository, article));
        }

        private async Task<bool> NotExists(IArticlesRepository repository, Article article)
        {
            return !await repository.Exists(article);
        }

        private ChannelReader<Article> StartProducing(IEnumerable<IArticlesParser> articlesParsers)
        {
            var channel = Channel.CreateUnbounded<Article>();
            var parsers = articlesParsers.Select(parser => Produce(parser, channel.Writer));

            Task.Run(async () =>
            {
                await Task.WhenAll(parsers);
                await channel.CompleteAsync();
            });
            return channel.Reader;
        }
        
        private static async Task Produce(IArticlesParser parser, ChannelWriter<Article> writer)
        {
            await foreach (var article in parser.Parse())
            {
                await writer.WriteAsync(article);
            }
        }
        
    }
}