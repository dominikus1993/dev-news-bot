using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Shared.Messages;

namespace DevNews.WebHooks.Application.Services
{
    public interface IWebHookNotifier
    {
        Task Notify(IList<Article> article);
    }
}