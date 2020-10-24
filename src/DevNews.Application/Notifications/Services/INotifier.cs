using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Akka.Types;

namespace DevNews.Application.Notifications.Services
{
    public interface INotifier
    {
        Task Notify(IEnumerable<Article> article);
    }
}