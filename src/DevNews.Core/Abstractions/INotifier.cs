using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using DevNews.Core.Model;

namespace DevNews.Core.Abstractions
{
    public interface INotifier
    {
        Task Notify(IEnumerable<Article> articles, CancellationToken cancellationToken = default);
    }

    public interface INotificationBroadcaster
    {
        Task Broadcast(IEnumerable<Article> articles, CancellationToken cancellationToken = default);
    }
}