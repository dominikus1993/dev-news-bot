using System;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Abstractions;
using DevNews.Core.Repository;

namespace DevNews.Core.UseCases
{
    public record ParseArticlesAndSendItParam(int ArticleQuantity);

    public class ParseArticlesAndSendItUseCase
    {
        private readonly IArticlesProvider _articlesProvider;
        private readonly IArticlesRepository _articlesRepository;
        private readonly INotificationBroadcaster _notificationBroadcaster;

        public ParseArticlesAndSendItUseCase(IArticlesProvider articlesProvider, IArticlesRepository articlesRepository,
            INotificationBroadcaster notificationBroadcaster)
        {
            _articlesProvider = articlesProvider;
            _articlesRepository = articlesRepository;
            _notificationBroadcaster = notificationBroadcaster;
        }

        public async Task Execute(ParseArticlesAndSendItParam param)
        {
            var articles = await _articlesProvider.Provide().OrderBy(_ => Guid.NewGuid()).Take(param.ArticleQuantity)
                .ToListAsync();
            await _articlesRepository.InsertMany(articles);
            await _notificationBroadcaster.Broadcast(articles);
        }
    }
}