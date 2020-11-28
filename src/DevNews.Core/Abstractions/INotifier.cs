using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Core.Model;

namespace DevNews.Core.Abstractions
{
    public interface INotifier
    {
        Task Notify(IEnumerable<Article> articles);
    }

    public interface INotificationBroadcaster
    {
        Task Broadcast(IEnumerable<Article> articles);
    }
}