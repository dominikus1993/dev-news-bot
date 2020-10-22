using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Shared.Messages;
using DevNews.Shared.Types;

namespace DevNews.WebHooks.Application.Services
{
    public interface IWebHookNotifier
    {
        Task Notify(IEnumerable<Article> article);
    }
}