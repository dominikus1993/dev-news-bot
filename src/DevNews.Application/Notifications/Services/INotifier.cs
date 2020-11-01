using System.Collections.Generic;
using System.Threading.Tasks;
using  DevNews.Core.Model;

namespace DevNews.Application.Notifications.Services
{
    public interface INotifier
    {
        Task Notify(IEnumerable<Article> article);
    }
}