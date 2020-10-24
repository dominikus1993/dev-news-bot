using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Shared.Types;
using DevNews.WebHooks.Application.Services;

namespace DevNews.WebHooks.Infrastructure.Services
{
    public class TeamsWebHookNotifier : INotifier
    {
        public async Task Notify(IEnumerable<Article> article)
        {
            throw new System.NotImplementedException();
        }
    }
}