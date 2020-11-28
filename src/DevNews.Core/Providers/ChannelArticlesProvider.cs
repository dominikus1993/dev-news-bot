using System.Collections.Generic;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using DevNews.Core.Repository;

namespace DevNews.Core.Providers
{
    public class ChannelArticlesProvider : IArticlesProvider
    {
        private IEnumerable<IArticlesParser> _articlesParsers;
        private IArticlesRepository _articlesRepository;

        public ChannelArticlesProvider(IEnumerable<IArticlesParser> articlesParsers, IArticlesRepository articlesRepository)
        {
            _articlesParsers = articlesParsers;
            _articlesRepository = articlesRepository;
        }

        public IAsyncEnumerable<Article> Provide()
        {
            throw new System.NotImplementedException();
        }
        
        
    }
}