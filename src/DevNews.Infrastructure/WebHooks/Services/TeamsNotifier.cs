using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Akka.Types;
using DevNews.Application.Notifications.Services;

namespace DevNews.Infrastructure.WebHooks.Services
{
    public class TeamsWebHookNotifier : INotifier
    {
        public async Task Notify(IEnumerable<Article> article)
        {
            throw new System.NotImplementedException();
        }
    }
}