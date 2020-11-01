using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Application.Notifications.Services;
using DevNews.Core.Model;

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