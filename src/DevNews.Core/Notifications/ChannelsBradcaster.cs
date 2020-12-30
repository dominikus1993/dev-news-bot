using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;

namespace DevNews.Core.Notifications
{
    public class ChannelsNotificationBroadcaster : INotificationBroadcaster
    {
        private readonly IEnumerable<INotifier> _notifiers;

        public ChannelsNotificationBroadcaster(IEnumerable<INotifier> notifiers)
        {
            _notifiers = notifiers;
        }

        public async Task Broadcast(IEnumerable<Article> articles)
        {
            var tasks = _notifiers.Select(notifier => notifier.Notify(articles));
            await Task.WhenAll(tasks);
        }
    }
}